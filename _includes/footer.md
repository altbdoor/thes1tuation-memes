{% comment %}
note: html layout files cannot render markdown directly
https://talk.jekyllrb.com/t/rendering-markdown-inside-an-html-include/4186
{% endcomment %}

{% capture footer_content %}

---

Generated with [Jekyll](https://jekyllrb.com/) at <time datetime="{{ 'now' | date_to_xmlschema }}">{{ 'now' | date_to_long_string }}</time>.

{% endcapture %}

{{ footer_content | markdownify }}

<ul class="list-inline">
  {% for item in site.data.credits %}
    <li class="list-inline-item">
      <a href="{{ item.link }}">{{ item.name }}</a>
    </li>
  {% endfor %}
</ul>
