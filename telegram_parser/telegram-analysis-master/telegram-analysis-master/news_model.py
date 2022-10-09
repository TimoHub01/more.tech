from copy import copy, Error
from datetime import datetime
from string import punctuation

import nltk
import numpy as np
import pandas as pd
import pymorphy2
from nltk.corpus import stopwords
from pymystem3 import Mystem
from sklearn.cluster import AffinityPropagation
from sklearn.feature_extraction.text import TfidfVectorizer
nltk.download('stopwords')

class NewsModel:

    def __init__(self, accountant_token_dict, random_state=5):
        self.accountant_token_dict = accountant_token_dict
        self.random_state = random_state
        self.mystem = Mystem()

        stopwords.words("russian").extend(['«', '="', '»'])
        self.russian_stopwords = stopwords.words("russian")

        self.tfidf = None
        self.clustering = None

    @staticmethod
    def pos(word, morth=pymorphy2.MorphAnalyzer()):
        return morth.parse(word)[0].tag.POS

    def preprocess_text(self, text):
        functors_pos = {'ADVB', 'NPRO', 'INTJ', 'PRED', 'PRCL', 'CONJ', 'PREP', 'PNCT', 'NUMB'}
        tokens = self.mystem.lemmatize(text.lower())
        tokens = [token for token in tokens if token not in self.russian_stopwords \
                  and token != " " \
                  and token.strip() not in punctuation + '«="»-,'
                  and self.pos(token) not in functors_pos
                  and all(map(lambda x: not x.isdigit(), token))]

        tokens = list(map(lambda x: x.replace('»', '').replace('«', '').replace('-', ''), tokens))

        text = " ".join(tokens)

        return text

    def score_word(self, row):
        result = 0
        for token in self.preprocess_text(row).split():
            try:
                result += self.accountant_token_dict[token]
            except:
                result += 0
        return result

    def markup_roles(self, text):
        target_score = 0.014518
        keywords = ['центробанк', 'цб', 'банк', 'экономика', 'правительство', 'вклад', 'оборот', 'государственный',
                    'минфин', 'акции', 'кредит', 'курс', 'биржа', 'рынок', 'сбербанк', 'втб']

        for word in text.split():
            if word in keywords:
                return 1

        result = self.score_word(text)
        if result >= target_score:
            return 2
        else:
            return 0

    def fit_predict(self, df, text_name='text', title_name='topic'):
        self.df = df
        self.df['role'] = self.df[text_name].apply(self.markup_roles)

        tfidf_vectorizer = TfidfVectorizer(preprocessor=self.preprocess_text)
        print('TF-IDF получен')
        df = self.df.loc[:, [text_name, title_name]].dropna()
        text = list(df[text_name].values)
        title = list(df[title_name].values)

        all_text = copy(text)
        #         all_text.extend(title)

        self.tfidf = tfidf_vectorizer.fit(text)

        text_tfidf = self.tfidf.transform(text)

        cluster_labels = AffinityPropagation(random_state=self.random_state).fit_predict(text_tfidf)
        print('Кластер обучен')

        self.clusters_count_ = dict(pd.Series(cluster_labels))

        self.df['text_cluster'] = pd.Series(cluster_labels)
        self.df['text_clust_cnt'] = self.df['text_cluster'].apply(lambda x: self.clusters_count_.get(x))

        group_data = self.df.groupby(['text_cluster']).agg({'views_cnt': 'sum', 'forwards_cnt': 'sum'}).reset_index()
        self.df = self.df.merge(group_data, how='left', on=['text_cluster'])
        self.df['views_cnt'] = self.df['views_cnt_y']
        self.df['forwards_cnt'] = self.df['forwards_cnt_y']
        self.df.drop(['views_cnt_y', 'views_cnt_x', 'forwards_cnt_y', 'forwards_cnt_x'], inplace=True)

        self.find_trend(title_name)
        print('Тренды найдены')

        self.relevant_score()
        print('Ранги расчитаны')

        return self.df

    def text_transform(self, df, text_name):
        if self.tfidf == None:
            raise Error('Не обучили модель!')

        return self.tfidf.transform(df[text_name].values)

    def predict(self, df, text_name):
        X = self.text_transform(df, text_name).todense()
        return self.clustering.predict(X)

    def find_trend(self, title_name):
        clusters = self.df.text_cluster.unique()
        self.df['trend'] = 0
        for i in clusters:
            if len(list(self.df[self.df.text_cluster == i].date.values)) > 0:
                min_date = datetime.timestamp(pd.to_datetime(min(self.df[self.df.text_cluster == i].date.values)))
                max_date = datetime.timestamp(pd.to_datetime(max(self.df[self.df.text_cluster == i].date.values)))
                month_sec = 60 * 60 * 24 * 30
                month_delta = (max_date - min_date) / month_sec
                if month_delta > 1:
                    pull_titles = self.df.loc[self.df.text_cluster == i, [title_name]].values.sum()
                    process_text = self.preprocess_text(pull_titles)
                    data_cons_split = process_text.split()
                    unique_words = list(set(data_cons_split))
                    occurrence = {}
                    for word in unique_words:
                        occurrence[word] = data_cons_split.count(word) / len(data_cons_split)
                    trend_name = max(occurrence, key=occurrence.get)
                    self.df.loc[self.df['text_cluster'] == i, ['trend']] = trend_name

    def relevant_score(self):
        self.df['relevant_score'] = 0

        def relevant_scoring(x):
            ttl = 1 / np.log(datetime.timestamp(datetime.now()) - datetime.timestamp(x['date']))
            if type(x['trend']) != str:
                trend = 0
            else:
                trend = 1
            x['relevant_score'] = ttl * 10000 + \
                                  x['text_clust_cnt'] * 1000 + \
                                  trend * 100 + \
                                  x['views_cnt'] * 10 + x['forwards_cnt']
            return x

        self.df = self.df.apply(relevant_scoring, axis=1)
        return self.df.sort_values(by=['relevant_score'])