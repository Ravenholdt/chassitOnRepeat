# Chassit on Repeat API

This is a documentation for the Chassit on Repeat API

## GET /api

<details>
    <summary><code>Returns a message that the api is running</code></summary>

##### Parameters

> None

##### Query

> None

##### Responses

> | http code | content-type               | response         |
> |-----------|----------------------------|------------------|
> | `200`     | `text/plain;charset=UTF-8` | `Api is running` |

##### Example cURL

> ```javascript
>  curl -X GET -H "Content-Type: application/json" http://localhost:8080/api
> ```

</details>


## GET /api/v1

<details>
    <summary><code>Returns a message that the api is running</code></summary>

##### Parameters

> None

##### Query

> None

##### Responses

> | http code | content-type               | response         |
> |-----------|----------------------------|------------------|
> | `200`     | `text/plain;charset=UTF-8` | `Api is running` |

##### Example cURL

> ```javascript
>  curl -X GET -H "Content-Type: application/json" http://localhost:8080/api/v1
> ```

</details>


## GET /api/v1/stats

<details>
    <summary><code>Returns playtime and video statistics</code></summary>

##### Parameters

> None

##### Query

> None

##### Responses

> | http code | content-type       | response                                                                                                 |
> |-----------|--------------------|----------------------------------------------------------------------------------------------------------|
> | `200`     | `application/json` | `{"repeat_entries":513,"total_time":18654866,"total_time_formatted":"215d 21h 54m 26s","video_count":2}` |

##### Example cURL

> ```javascript
>  curl -X GET -H "Content-Type: application/json" http://localhost:8080/api/v1/stats
> ```

</details>


## GET /api/v1/video

<details>
    <summary><code>Returns an array of all videos</code></summary>

##### Parameters

> None

##### Query

> None

##### Responses

> | http code | content-type       | response                                                                                                                                                                                                |
> |-----------|--------------------|---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
> | `200`     | `application/json` | Json Array with video info <br/>`[{"end":205,"id":"iR-K2rUP86M","safe":true,"start":0,"time":733080,"time_formatted":"8d 11h 38m 00s","title":"DAGOTHWAVE","url":"/files/DAGOTHWAVE-iR-K2rUP86M.mp4"}]` |

##### Example cURL

> ```javascript
>  curl -X GET -H "Content-Type: application/json" http://localhost:8080/api/v1/video
> ```

</details>


## GET /api/v1/video/random

<details>
    <summary><code>Returns a random video</code></summary>

##### Parameters

> None

##### Query

> | name   | type     | data type | description                            |
> |--------|----------|-----------|----------------------------------------|
> | `safe` | optional | flag      | If only safe videos should be returned |

##### Responses

> | http code | content-type       | response                                                                                                                                                              |
> |-----------|--------------------|-----------------------------------------------------------------------------------------------------------------------------------------------------------------------|
> | `200`     | `application/json` | `{"end":205,"id":"iR-K2rUP86M","safe":true,"start":0,"time":733080,"time_formatted":"8d 11h 38m 00s","title":"DAGOTHWAVE","url":"/files/DAGOTHWAVE-iR-K2rUP86M.mp4"}` |
> | `404`     | `application/json` | `{"error":{"code":403,"message":"Error getting random video"},"message":"There is nothing here...","path":"/api/v1/video/random"}`                                    |


##### Example cURL

> ```javascript
>  curl -X GET -H "Content-Type: application/json" http://localhost:8080/api/v1/video/random
> ```


> ```javascript
>  curl -X GET -H "Content-Type: application/json" http://localhost:8080/api/v1/video/random?safe
> ```

</details>


## GET /api/v1/video/{id}

<details>
    <summary><code>Returns a specific video</code></summary>

##### Parameters

> | name |  type      | data type      | description                           |
> |------|------------|----------------|---------------------------------------|
> | `id` |  required  | string         | The id of the video to get info about |

##### Query

> None

##### Responses

> | http code | content-type       | response                                                                                                                                                              |
> |-----------|--------------------|-----------------------------------------------------------------------------------------------------------------------------------------------------------------------|
> | `200`     | `application/json` | `{"end":205,"id":"iR-K2rUP86M","safe":true,"start":0,"time":733080,"time_formatted":"8d 11h 38m 00s","title":"DAGOTHWAVE","url":"/files/DAGOTHWAVE-iR-K2rUP86M.mp4"}` |
> | `404`     | `application/json` | `{"error":{"code":403,"message":"Video not found"},"message":"There is nothing here...","path":"/api/v1/video/iR-K2rUP86M"}`                                          |

##### Example cURL

> ```javascript
>  curl -X GET -H "Content-Type: application/json" http://localhost:8080/api/v1/video/iR-K2rUP86M
> ```

</details>


## POST /api/v1/video/{id}

<details>
    <summary><code>Updates the repeated time of the specified video</code></summary>

##### Parameters

> | name |  type      | data type      | description                   |
> |------|------------|----------------|-------------------------------|
> | `id` |  required  | string         | The id of the video to update |

##### Post data
> | name   |  type      | data type | description                           |
> |--------|------------|-----------|---------------------------------------|
> | `time` |  required  | int       | The repeated time to add to the video |

##### Responses

> | http code | content-type       | response                                                                                                                                                                                         |
> |-----------|--------------------|--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
> | `200`     | `application/json` | The updated video info<br/>`{"end":205,"id":"iR-K2rUP86M","safe":true,"start":0,"time":733090,"time_formatted":"8d 11h 38m 10s","title":"DAGOTHWAVE","url":"/files/DAGOTHWAVE-iR-K2rUP86M.mp4"}` |
> | `400`     | `application/json` | `{"error":{"code":400,"message":"Bad body"},"message":"There is nothing here...","path":"/api/v1/video/iR-K2rUP86M"}`                                                                            |

##### Example cURL

> ```javascript
>  curl -X POST -H "Content-Type: application/json" -d '{"time": 10}' http://localhost:8080/api/v1/video/iR-K2rUP86M
> ```

</details>

## POST /api/v1/video/{id}/settings

<details>
    <summary><code>Updates the start/end and safe status of the specified video</code></summary>

##### Parameters

> | name |  type      | data type      | description                   |
> |------|------------|----------------|-------------------------------|
> | `id` |  required  | string         | The id of the video to update |

##### Post data
> | name    | type     | data type | description                        |
> |---------|----------|-----------|------------------------------------|
> | `safe`  | required | boolean   | If the video is safe to the public |
> | `start` | optional | int       | The start of the repeat in seconds |
> | `end`   | optional | int       | The end of the repeat in seconds   |

##### Responses

> | http code | content-type       | response                                                                                                                                                                                         |
> |-----------|--------------------|--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
> | `200`     | `application/json` | The updated video info<br/>`{"end":205,"id":"iR-K2rUP86M","safe":true,"start":0,"time":733090,"time_formatted":"8d 11h 38m 10s","title":"DAGOTHWAVE","url":"/files/DAGOTHWAVE-iR-K2rUP86M.mp4"}` |
> | `400`     | `application/json` | `{"error":{"code":400,"message":"Bad body"},"message":"There is nothing here...","path":"/api/v1/video/iR-K2rUP86M/settings"}`                                                                   |

##### Example cURL

> ```javascript
>  curl -X POST -H "Content-Type: application/json" -d '{"start": null, "end": 200, "safe": true}' http://localhost:8080/api/v1/video/iR-K2rUP86M/settings
> ```

</details>

## GET /api/v1/playlist

<details>
    <summary><code>Returns an array of all playlists</code></summary>

##### Parameters

> None

##### Query

> None

##### Responses

> | http code | content-type       | response                                                                                                                                                                                                                                                                                            |
> |-----------|--------------------|-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
> | `200`     | `application/json` | Json Array with playlist info <br/>`[{"created_at":"2023-11-13T18:15:00Z","updated_at":"2023-11-15T22:42:39.075Z","id":"hMNXaW","name":"Unsafe Playlist","playtime":1815,"last_played":"2023-11-15T22:42:39.075Z","videos":["Y9EKzvTo3g0","gDtElkw_o6w"],"safe":false,"time_formatted":"30m 15s"}]` |

##### Example cURL

> ```javascript
>  curl -X GET -H "Content-Type: application/json" http://localhost:8080/api/v1/playlist
> ```

</details>


## GET /api/v1/playlist/{id}

<details>
    <summary><code>Returns a specific playlist</code></summary>

##### Parameters

> | name |  type      | data type      | description                              |
> |------|------------|----------------|------------------------------------------|
> | `id` |  required  | string         | The id of the playlist to get info about |

##### Query

> None

##### Responses

> | http code | content-type       | response                                                                                                                                                                                                                                                       |
> |-----------|--------------------|----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
> | `200`     | `application/json` | `{"created_at":"2023-11-13T18:15:00Z","updated_at":"2023-11-15T22:42:39.075Z","id":"hMNXaW","name":"Unsafe Playlist","playtime":1815,"last_played":"2023-11-15T22:42:39.075Z","videos":["Y9EKzvTo3g0","gDtElkw_o6w"],"safe":false,"time_formatted":"30m 15s"}` |
> | `404`     | `application/json` | `{"error":{"code":403,"message":"Playlist not found"},"message":"There is nothing here...","path":"/api/v1/playlist/hMNXaW"}`                                                                                                                                  |

##### Example cURL

> ```javascript
>  curl -X GET -H "Content-Type: application/json" http://localhost:8080/api/v1/playlist/hMNXaW
> ```

</details>


## POST /api/v1/playlist/{id}

<details>
    <summary><code>Updates the repeated time of the specified playlist</code></summary>

##### Parameters

> | name |  type      | data type      | description                      |
> |------|------------|----------------|----------------------------------|
> | `id` |  required  | string         | The id of the playlist to update |

##### Post data
> | name   |  type      | data type | description                              |
> |--------|------------|-----------|------------------------------------------|
> | `time` |  required  | int       | The repeated time to add to the playlist |

##### Responses

> | http code | content-type       | response                                                                                                                                                                                                                                                                                     |
> |-----------|--------------------|----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
> | `200`     | `application/json` | The updated playlist info<br/>`{"created_at":"2023-11-13T18:15:00Z","updated_at":"2023-11-15T22:42:39.075Z","id":"hMNXaW","name":"Unsafe Playlist","playtime":1815,"last_played":"2023-11-15T22:42:39.075Z","videos":["Y9EKzvTo3g0","gDtElkw_o6w"],"safe":false,"time_formatted":"30m 15s"}` |
> | `400`     | `application/json` | `{"error":{"code":400,"message":"Bad body"},"message":"There is nothing here...","path":"/api/v1/video/hMNXaW"}`                                                                                                                                                                             |

##### Example cURL

> ```javascript
>  curl -X POST -H "Content-Type: application/json" -d '{"time": 10}' http://localhost:8080/api/v1/playlist/hMNXaW
> ```

</details>


## GET /api/v1/playlist/{id}/random

<details>
    <summary><code>Returns a random video from a playlist</code></summary>

##### Parameters

> None

##### Query

> None

##### Responses

> | http code | content-type       | response                                                                                                                                                              |
> |-----------|--------------------|-----------------------------------------------------------------------------------------------------------------------------------------------------------------------|
> | `200`     | `application/json` | `{"end":205,"id":"iR-K2rUP86M","safe":true,"start":0,"time":733080,"time_formatted":"8d 11h 38m 00s","title":"DAGOTHWAVE","url":"/files/DAGOTHWAVE-iR-K2rUP86M.mp4"}` |
> | `404`     | `application/json` | `{"error":{"code":403,"message":"Error getting random video"},"message":"There is nothing here...","path":"/api/v1/playlist/hMNXaW/random"}`                          |


##### Example cURL

> ```javascript
>  curl -X GET -H "Content-Type: application/json" http://localhost:8080/api/v1/playlist/hMNXaW/random
> ```

</details>