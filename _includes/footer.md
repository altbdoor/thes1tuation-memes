{% capture footer_content %}

---

Generated with [Jekyll](https://jekyllrb.com/) at {{ 'now' | date_to_long_string }}.

{% endcapture %}

{{ footer_content | markdownify }}

<ul class="list-inline">
  {% for item in site.data.credits %}
  <li class="list-inline-item">
    <a href="{{ item.link }}">{{ item.name }}</a>
  </li>
  {% endfor %}
</ul>
