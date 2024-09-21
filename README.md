# vitess-plugin-poc

```bash
tim@Tims-MacBook-Pro vitess-plugin-poc % make
go build -buildmode=plugin ./durabler/cross_cell.go
go build -o plugin-poc ./main.go ./durability.go
tim@Tims-MacBook-Pro vitess-plugin-poc % ls -alh cross_cell.so plugin-poc 
-rw-r--r--  1 tim  staff   6.7M Sep 21 00:42 cross_cell.so
-rwxr-xr-x  1 tim  staff   8.3M Sep 21 00:42 plugin-poc
tim@Tims-MacBook-Pro vitess-plugin-poc % ./plugin-poc -plugin-path cross_cell.so
plugin path: cross_cell.so
plugin: &{plugin/unnamed-5dfeccb51fb22fa97e7ff9a6985e5883f57f03b5  0xc00008c0c0 map[DurabilityCrossCell:0x150f64bd0]}
must_not
```
