# v2.0.0

* **Breaking:** `WriteConfErr` now requires an `error` argument (or `nil`). When an error is provided, its message is sent to ThingsDB; pass `nil` to preserve the previous behavior of sending an empty package.
