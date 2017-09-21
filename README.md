# gproxy

## check proxy (chrome)
in console log
```js
fetch("http://localhost:8080/search/repositories?q=golang")
    .then(data => {
         data.json()
         .then(y => console.log(y));
     });
```