#!/usr/bin/env python

import os
import json
from typing import TypedDict
from itertools import groupby
from datetime import datetime

current_dir = os.path.dirname(__file__)
raw_data = os.path.join(current_dir, "../data/discord.json")
parsed_data = os.path.join(current_dir, "../data/discord-parsed.json")

Post = TypedDict(
    "Post",
    {
        "filename": str,
        "timestamp": str,
        "tags": list[str],
    }
)

post_list: list[Post] = []

with open(raw_data, 'r') as fp:
    post_list = json.load(fp)

for datum in post_list:
    post_datetime = datetime.strptime(datum["timestamp"], "%m/%d/%Y %I:%M %p")
    datum["epoch"] = post_datetime.timestamp()
    datum["timeDisplay"] = post_datetime.strftime("%d %b, %I:%M %p")
    datum["groupBy"] = post_datetime.strftime("%B %Y")

with open(parsed_data, "w") as fp:
    post_list = sorted(post_list, key=lambda x: x["epoch"], reverse=True)
    groups = groupby(post_list, key=lambda x: x["groupBy"])

    group_data = [{"name": key, "items": list(items)} for key, items in groups]
    json.dump(group_data, fp, indent=4)
