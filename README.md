# Go Wiki
## An example wiki written in go
wiki created by following [this](https://go.dev/doc/articles/wiki/ "https://go.dev/doc/articles/wiki/") example.

This branch contains the implementations the other features mentioned in the above article.

###### Todo

* Spruce up the page templates by making them valid HTML and adding some CSS rules.
* Implement inter-page linking by converting instances of [PageName] to
 <a href="/view/PageName">PageName</a>. (hint: you could use regexp.ReplaceAllFunc to do this
)

###### Implemented
* Store templates in tmpl/ and page data in data/.
* Add a handler to make the web root redirect to /view/FrontPage.
