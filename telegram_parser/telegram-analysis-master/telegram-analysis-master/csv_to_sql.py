import pandas as pd
import psycopg2 as psycopg2
import sqlalchemy
from sqlalchemy import create_engine


def csv_to_sql():
    data = pd.read_csv('TEST.csv')
    df = pd.DataFrame(data)

    conn = psycopg2.connect(
        host="localhost",
        database="postgres",
        user="admin",
        password="more-tech",
        port="15432")
    conn.set_client_encoding('utf-8')
    cursor = conn.cursor()


    # cursor.execute('''CREATE TABLE news (news_id serial primary key NOT NULL , date_news date, message varchar(5000),views float(50),
    # forwards float(50),title varchar)''')
    # conn.commit()
    #
    # cursor.execute('''CREATE TABLE news_new (id serial primary key NOT NULL , link varchar, topic varchar,
    # text varchar, date date, views_cnt float, forwsrd_cnt float, reactions_cnt integer)''')
    # conn.commit()

    # date_news, message, views, forwards, title
    # Insert DataFrame to Table
    for row in df.itertuples():

        print(row.date)
        print(row.date[0:11])
        cursor.execute("INSERT INTO news (date_news, message, views, forwards, title) VALUES (%s, %s, %s, %s, %s)", (row.date[0:11], row.message, row.views, row.forwards, row.title))

    conn.commit()
    conn.close()

def to_db(dataFrame):
    conn = psycopg2.connect(
        host="localhost",
        database="postgres",
        user="admin",
        password="more-tech",
        port="15432")
    cursor = conn.cursor()
    conn.set_client_encoding('utf-8')
    for row in dataFrame.itertuples():
        print(row.date)
        print(row.date[0:11])
        (print("new"))
        cursor.execute("INSERT INTO news_new (link, topic, text, date, views_cnt, forwards_cnt, reactions_cnt) VALUES (%s, %s, %s, %s, %s, %s, %s)", (row.link, row.topic, row.text, row.date, row.views_cnt,
                                                                                                                                                   row.forwards_cnt , row.reactions_cnt))
    df = cursor.fetchall()
    conn.commit()
    conn.close()



def get_data_for_model():
    conn = psycopg2.connect(
        host="localhost",
        database="postgres",
        user="admin",
        password="more-tech",
        port="15432")
    cursor = conn.cursor()
    conn.set_client_encoding('utf-8')
    # cursor.execute("SELECT date_news, message, views, forwards, title FROM news WHERE date_news BETWEEN '2022-01-01' AND '2022-12-31'")
    engine = create_engine("postgresql+psycopg2://admin:more-tech@localhost:15432/postgres?client_encoding=utf8")
    with engine.connect() as con:
        df = pd.read_sql("SELECT date_news as date, message as text, views as views_cnt, forwards as forwards_cnt, title as topic FROM news WHERE date_news BETWEEN '2022-01-01' AND '2022-12-31'", con)
    print(df)
    # df.to_sql('news', connect)

    # df.loc[:,["date"]] = df.loc[:,["date_news"]]
    # # df['date_news'] = df['date']
    # df.drop(['date_news'], axis=1, inplace=True)
    # df['text'] = df['message']
    # df.drop(['message'], axis=1, inplace=True)
    # df['views_cnt'] = df['views']
    # df.drop(['views'], axis=1, inplace=True)
    # df['forwards_cnt'] = df['forwards']
    # df.drop(['forwards'], axis=1, inplace=True)
    # df['topic'] = df['title']
    # df.drop(['title'], axis=1, inplace=True)
    df['link'] = 'null'
    df['reactions_cnt'] = 0

    conn.commit()
    conn.close()

    return df

