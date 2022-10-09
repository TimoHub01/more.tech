from typing import Optional, List, Any
from datetime import datetime


class Entity:
    empty: str
    offset: int
    length: int
    url: Optional[str]

    def __init__(self, empty: str, offset: int, length: int, url: Optional[str]) -> None:
        self.empty = empty
        self.offset = offset
        self.length = length
        self.url = url

class PeerID:
    empty: str
    channel_id: int

    def __init__(self, empty: str, channel_id: int) -> None:
        self.empty = empty
        self.channel_id = channel_id


class Result:
    empty: str
    reaction: str
    count: int
    chosen: bool

    def __init__(self, empty: str, reaction: str, count: int, chosen: bool) -> None:
        self.empty = empty
        self.reaction = reaction
        self.count = count
        self.chosen = chosen


class Reactions:
    empty: str
    results: List[Result]
    min: bool
    can_see_list: bool
    recent_reactions: List[Any]

    def __init__(self, empty: str, results: List[Result], min: bool, can_see_list: bool, recent_reactions: List[Any]) -> None:
        self.empty = empty
        self.results = results
        self.min = min
        self.can_see_list = can_see_list
        self.recent_reactions = recent_reactions


class Welcome4:
    empty: str
    id: int
    peer_id: PeerID
    date: datetime
    message: str
    out: bool
    mentioned: bool
    media_unread: bool
    silent: bool
    post: bool
    from_scheduled: bool
    legacy: bool
    edit_hide: bool
    pinned: bool
    noforwards: bool
    from_id: None
    fwd_from: None
    via_bot_id: None
    reply_to: None
    media: None
    reply_markup: None
    entities: List[Entity]
    views: int
    forwards: int
    replies: None
    edit_date: datetime
    post_author: None
    grouped_id: None
    reactions: Reactions
    restriction_reason: List[Any]
    ttl_period: None

    def __init__(self, empty: str, id: int, peer_id: PeerID, date: datetime, message: str, out: bool, mentioned: bool, media_unread: bool, silent: bool, post: bool, from_scheduled: bool, legacy: bool, edit_hide: bool, pinned: bool, noforwards: bool, from_id: None, fwd_from: None, via_bot_id: None, reply_to: None, media: None, reply_markup: None, entities: List[Entity], views: int, forwards: int, replies: None, edit_date: datetime, post_author: None, grouped_id: None, reactions: Reactions, restriction_reason: List[Any], ttl_period: None) -> None:
        self.empty = empty
        self.id = id
        self.peer_id = peer_id
        self.date = date
        self.message = message
        self.out = out
        self.mentioned = mentioned
        self.media_unread = media_unread
        self.silent = silent
        self.post = post
        self.from_scheduled = from_scheduled
        self.legacy = legacy
        self.edit_hide = edit_hide
        self.pinned = pinned
        self.noforwards = noforwards
        self.from_id = from_id
        self.fwd_from = fwd_from
        self.via_bot_id = via_bot_id
        self.reply_to = reply_to
        self.media = media
        self.reply_markup = reply_markup
        self.entities = entities
        self.views = views
        self.forwards = forwards
        self.replies = replies
        self.edit_date = edit_date
        self.post_author = post_author
        self.grouped_id = grouped_id
        self.reactions = reactions
        self.restriction_reason = restriction_reason
        self.ttl_period = ttl_period
