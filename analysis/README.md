# mercari analysis

This package doesn't contain any function that called by other packages, all 
of funcs in this package are test functions that help to analyze how mercari 
api v2 works (for test functions are easy to run in vs code).

## Suspicious Elliptic Curve Key Pair

Appears in jp.mercari.com's Local Storage.

Keys are self-generated.

Private Key:

```json
{
  "crv": "P-256",
  "d": "kRAB_FrWLQdSWnCiBFnjrhG6Jbqg6Fwtm7QrDYPC0mg",
  "ext": true,
  "key_ops": [
    "sign"
  ],
  "kty": "EC",
  "x": "e5zwA7scbeI2653Z6hKV-ktJ9fDXAIce8GLDWcCl-Z0",
  "y": "9QXtjg2OCbLmesYFDl50u7dI690wMINClVz2HgZQNng"
}
```

Public Key:

```json
{
  "crv": "P-256",
  "ext": true,
  "key_ops": [
    "verify"
  ],
  "kty": "EC",
  "x": "e5zwA7scbeI2653Z6hKV-ktJ9fDXAIce8GLDWcCl-Z0",
  "y": "9QXtjg2OCbLmesYFDl50u7dI690wMINClVz2HgZQNng"
}
```

Use Test_buildECPrivKeys() to parse key to PEM format.

## JWT

Part 1:

```json
{
    "typ": "dpop+jwt",
    "alg": "ES256",
    "jwk": {
        "crv": "P-256",
        "kty": "EC",
        "x": "e5zwA7scbeI2653Z6hKV-ktJ9fDXAIce8GLDWcCl-Z0",
        "y": "9QXtjg2OCbLmesYFDl50u7dI690wMINClVz2HgZQNng"
    }
}
```

A perfect match with keys in Local Storage. Using SHA256 as hash.

Part 2:

```json
{
    "iat": 1689558892,
    "jti": "870ccfcd-aad0-469d-a3ac-a5a88adfd9af",
    "htu": "https://api.mercari.jp/v2/entities:search",
    "htm": "POST",
    "uuid": "f25aeb43-513d-4045-94a8-23fa4618b265"
}
```

uuid is generated randomly.

Part 3:
```
J4EPhmNia_4AQfKenUi8xtSV94ru9DpXesx-1F-mh5-q1zYNpSOYvR7d7ERl9OcZGFDj9PYu51UzBSmVVhhPgA
```

Signature is not asn1 encoded. Just concat r and s and get 64 bytes.

## How searchSessionId generated

At webpack 76152 (line 86714) (might change in future):

```js
var f = 32;

function h() {
    var a = arguments.length > 0 && void 0 !== arguments[0] ? arguments[0] : f
      , b = Math.ceil(a / 2);
    return g(e.e.getRandomValues(new Uint8Array(b))).slice(0, a)
}
```

Full:

```js
76152: function(a, b, c) {
  "use strict";
  c.d(b, {
      O: function() {
          return h
      }
  });
  var d = c(69779)
    , e = c(91371)
    , f = 32
    , g = function(a) {
      return (0,
      d.Z)(a).map(function(a) {
          return a.toString(16).padStart(2, "0")
      }).join("")
  };
  function h() {
      var a = arguments.length > 0 && void 0 !== arguments[0] ? arguments[0] : f
        , b = Math.ceil(a / 2);
      return g(e.e.getRandomValues(new Uint8Array(b))).slice(0, a)
  }
},
```

Rewrite:
```js
function generateSearchSessionId() {
    var a = arguments.length;
    if (a > 0 && void 0 !== arguments[0]) {
      a = arguments[0];
    } else {
      a = 32;
    }
    length = Math.ceil(a / 2); // Uint8Array length
    return AnyFunctionThatConvertByteToLowerCaseHexString(window.crypto.getRandomValues(new Uint8Array(length))).slice(0, a)
}
```

It's the second time that function called when entities:search generating arguments.