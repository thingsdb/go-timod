# v2.0.1

* Added `WriteResponseNil` for simple `nil` response.
* Added `NilBytes` constant, representing `nil` bytes in MessagePack format.

# v2.0.0

* **Breaking:** `WriteConfErr` now requires an `error` argument (or `nil`). When an error is provided, its message is sent to ThingsDB; pass `nil` to preserve the previous behavior of sending an empty package.
