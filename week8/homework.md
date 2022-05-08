## 1. 使用 redis benchmark 工具, 测试 10 20 50 100 200 1k 5k 字节 value 大小，redis get set 性能。

测试脚本

```bash
#!/usr/bin/bash
n=${1:-100000}
for SIZE in 10 20 50 100 200 1000 2000 5000
do 
  echo redis-benchmark -q -d $SIZE -t get,set -n $n
  redis-benchmark -q -d $SIZE -t get,set -n $n
done
```

100000个请求

```
redis-benchmark -q -d 10 -t get,set -n 100000
SET: 164473.69 requests per second, p50=0.151 msec                    
GET: 156739.81 requests per second, p50=0.159 msec                    

redis-benchmark -q -d 20 -t get,set -n 100000
SET: 167224.08 requests per second, p50=0.151 msec                    
GET: 158478.61 requests per second, p50=0.159 msec                    

redis-benchmark -q -d 50 -t get,set -n 100000
SET: 168067.22 requests per second, p50=0.151 msec                    
GET: 166944.92 requests per second, p50=0.151 msec                    

redis-benchmark -q -d 100 -t get,set -n 100000
SET: 161030.59 requests per second, p50=0.159 msec                    
GET: 167224.08 requests per second, p50=0.151 msec                    

redis-benchmark -q -d 200 -t get,set -n 100000
SET: 157728.70 requests per second, p50=0.167 msec                    
GET: 179533.22 requests per second, p50=0.143 msec                    

redis-benchmark -q -d 1000 -t get,set -n 100000
SET: 159744.41 requests per second, p50=0.159 msec                    
GET: 181818.17 requests per second, p50=0.135 msec                    

redis-benchmark -q -d 2000 -t get,set -n 100000
SET: 157232.70 requests per second, p50=0.167 msec                    
GET: 179211.45 requests per second, p50=0.143 msec                    

redis-benchmark -q -d 5000 -t get,set -n 100000
SET: 152207.00 requests per second, p50=0.167 msec                    
GET: 169204.73 requests per second, p50=0.151 msec
```

## 2. 写入一定量的 kv 数据, 根据数据大小 1w-50w 自己评估, 结合写入前后的 info memory 信息 , 分析上述不同 value 大小下，平均每个 key 的占用内存空间。

写入前

```bash
$ redis-cli flushall
$ redis-cli info memory
# Memory
used_memory:1355624
used_memory_human:1.29M
used_memory_rss:11837440
used_memory_rss_human:11.29M
used_memory_peak:5572656
used_memory_peak_human:5.31M
used_memory_peak_perc:24.33%
used_memory_overhead:864752
used_memory_startup:862768
used_memory_dataset:490872
used_memory_dataset_perc:99.60%
allocator_allocated:1459952
allocator_active:2830336
allocator_resident:7405568
total_system_memory:67348508672
total_system_memory_human:62.72G
used_memory_lua:31744
used_memory_vm_eval:31744
used_memory_lua_human:31.00K
used_memory_scripts_eval:0
number_of_cached_scripts:0
number_of_functions:0
number_of_libraries:0
used_memory_vm_functions:32768
used_memory_vm_total:64512
used_memory_vm_total_human:63.00K
used_memory_functions:184
used_memory_scripts:184
used_memory_scripts_human:184B
maxmemory:0
maxmemory_human:0B
maxmemory_policy:noeviction
allocator_frag_ratio:1.94
allocator_frag_bytes:1370384
allocator_rss_ratio:2.62
allocator_rss_bytes:4575232
rss_overhead_ratio:1.60
rss_overhead_bytes:4431872
mem_fragmentation_ratio:8.87
mem_fragmentation_bytes:10502488
mem_not_counted_for_evict:0
mem_replication_backlog:0
mem_total_replication_buffers:0
mem_clients_slaves:0
mem_clients_normal:1800
mem_cluster_links:0
mem_aof_buffer:0
mem_allocator:jemalloc-5.2.1
active_defrag_running:0
lazyfree_pending_objects:0
lazyfreed_objects:0
```

写入数据

```bash
#!/usr/bin/bash
size=${1:-10}
for ((i=0;i<10000;i++))
do
    echo -en python3 -c 'print("A" * $size)' | redis-cli -x set key$i >>redis.log
done
```

1. 写入10000条数据, 其中 key长度为4，value长度为10

```bash
$ redis-cli info memory
# Memory
used_memory:2214032
used_memory_human:2.11M

$ redis-cli memory usage key1
(integer) 64
```

2. 写入10000条数据, 其中 key长度为4，value长度为20

```bash
# Memory
used_memory:2454032
used_memory_human:2.34M

$ redis-cli memory usage key1
(integer) 72
```

3 写入10000条数据, 其中 key长度为4，value长度为50

```bash
# Memory
used_memory:2614032
used_memory_human:2.49M

$ redis-cli memory usage key1
(integer) 104
```

4.写入10000条数据, 其中 key长度为4，value长度为100

```bash
# Memory
used_memory:3174032
used_memory_human:3.03M

$ redis-cli memory usage key1
(integer) 160
```

5.写入10000条数据, 其中 key长度为4，value长度为200

```bash
# Memory
used_memory:4318688
used_memory_human:4.12M

$ redis-cli memory usage key1
(integer) 272
```

6.写入10000条数据, 其中 key长度为4，value长度为1000

```bash
# Memory
used_memory:12318688
used_memory_human:11.75M

$ redis-cli memory usage key1
(integer) 1072
```

7.写入10000条数据, 其中 key长度为4，value长度为2000

```bash
~ ❯ redis-cli info memory      
# Memory
used_memory:22558688
used_memory_human:21.51M

$ redis-cli memory usage key1
(integer) 2096
```

8.写入10000条数据, 其中 key长度为4，value长度为5000

```bash
# Memory
used_memory:53278688
used_memory_human:50.81M

$ redis-cli memory usage key1
(integer) 5168
```