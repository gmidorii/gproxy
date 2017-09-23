# gproxy

### API Doc
#### Request Parameter
| Name      | Type   | Default | Description                         |
|-----------+--------+---------+-------------------------------------|
| cors-host | string |         | The host name which want to request |
| proto     | string | http    | The protocol connection             |


### check proxy (chrome)
in console log
```js
fetch("http://localhost:8080/search/repositories?q=golang")
    .then(data => {
         data.json()
         .then(y => console.log(y));
     });
```
