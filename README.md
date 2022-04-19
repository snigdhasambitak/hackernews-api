# hackernews-api
An HTTP API micro-service that uses the YCombinator HackerNews API and upon request returns the top 50 of the latest 500 stories where the author has karma above 2413. The position is determined by the number of comments in relation to the top 50 stories. The story with the most comments should have position: 1, so on and so forth.

# Specification

## Request:

GET /stories

## Response:

```json
{
"stories": [
{
"author": "nick1",
"karma": 5341,
"comments": 192,
"title": "article title", "position": 1
},
{
"author": "nick2",
"karma": 7629,
"title": "article title", "comments": 12,
"position": 3
},
{
"author": "nick3",
"karma": 6293,
"title": "article title", "comments": 180,
"position": 2
}
]
}
```

# Implementation

As the ask is to get a curated list of the top 50 stories based on the authors karma, so we need to first set up different models for our requirements.
We will break down the application into 3 different parts

1. `Models` : We define the structs for the various items. The models which we have used are :
   1. `items`(Stories, comments, jobs, Ask HNs and even polls are just items. They're identified by their ids, which are unique integers, and live under /v0/item/<id>), 
   2. `story`(the response structure containing Author, Karma, Title, Comments and Position), 
   3. `topstories`(the list of top stories) and 
   4. `user`(which contains the varius fileds of an author)
2. `Handlers` : Where we mock and call the various services. This also enables us to extend our program and later add additional functionalities instead of modifying the existing services
3. `Hackernews services` : We define the various services as per our requirememts.
   1. GetItem returns item from hackernews API for given id
      ```
      GetItem(id int) (models.Item, error)
      ```
   2. GetUser returns user from hackernews API for given username
      ```
      GetUser(username string) (models.User, error)
      ```
   3. GetTopStories returns top 500 stories from hackernews API
      ```
      GetTopStories() ([]models.Item, error)
      ```
   4. Curated50 returns top 50 of the latest 500 stories where the author has karma above 2413 with most comments
      ```
      Curated50(minKarma int) ([]models.Story, error)
      ```

* We first collect the curated50 stories based on the authors minimum karma i.e 2413 and then sort the list based on the comments. 

* We use goroutines( 50 workers ) for parallel execution of our code as it takes around 3 mins to go through the entire 500 list if we are using a serialisation method.

* We have also used a in memory caching mechanism that creates a cache with a default expiration time of 5 minutes, and which purges expired items every 10 minutes

```go
cache:      cache.New(5*time.Minute, 10*time.Minute),
```

# Deployment 

* `chart` folder contains helm chart files.
   * `hackernews-api`
      * `templates` folder contains the following:
         * `deployment.yaml` contains K8s deployment specifics and env level variables.
         * `hpa.yaml` is used to specify auto-pod scaling in the cluster
         * `NetworkPolicy.yaml` file manages external access to the services in a cluster, typically HTTP.
         * `pdb.yaml` file is used to manage Pod Distribution Budget that the cluster will honor
         * `service.yaml` specifies service management like the loadbalancer with AWS

     *`chart.yaml` is used to specify app-pod level information that is used throughout the K8s config

     *`values.yaml` is used to handle app information and specifics. we can use `values/dev.yaml`
     to override these values and use them for each environment.

# Local Development

The prerequisite is to have a local cluster based on minikube.

```shell
brew install minikube
minikube start
```

We can leverage the below makefile for building and deploying out application.

```makefile
OSTYPE ?= darwin
env ?= dev
app_name ?= hackernews-api
repo_root = ${PWD}
cluster_name ?= aws_test_cluster
namespace_name ?= hackernews-api
helm_root ?= ${PWD}/chart/hackernews-api
docker_repo ?= snigdhasambit/hackernews-api
docker_release ?= 1.0
KUBECONFIG ?= /root/.kube/config

set_namespace:
	kubectl config use-context ${cluster_name} --kubeconfig=${KUBECONFIG} \
  	kubectl config set-context ${cluster_name} --namespace ${namespace_name} --kubeconfig=${KUBECONFIG}

docker_build:
	docker build ${repo_root} -t ${docker_repo}:${docker_release}

docker_release: docker_build
	docker push ${docker_repo}:${docker_release}

deploy_dry:
	helm upgrade -i ${app_name} ${helm_root} \
    --set ImageVersion=${docker_release} \
    --debug \
    --dry-run

deploy: 
	helm upgrade -i ${app_name} ${helm_root} \
	--set ImageVersion=${docker_release} \
    --debug

service:
	minikube service hackernews-api --url

destroy:
	helm delete ${app_name}

# optional installation of prometheus

prometheus_deploy:
	helm repo add prometheus-community https://prometheus-community.github.io/helm-charts \
	helm install prometheus prometheus-community/prometheus \
    kubectl expose service prometheus-server --type=NodePort --target-port=9090 --name=prometheus-server-np

prometheus_url:
	minikube service prometheus-server-np --url

destroy_prometheus:
	helm delete prometheus


```

# Service endpoints

1. `GET /stories`

```json
[
{
"author": "ibobev",
"karma": 4135,
"comments": 276,
"title": "Unreal vs. Unity Opinion",
"position": 10
},
{
"author": "djoldman",
"karma": 3880,
"comments": 78,
"title": "What is the origin of “daemon” with regards to computing? (2011)",
"position": 50
},
{
"author": "cube00",
"karma": 2552,
"comments": 92,
"title": "When you get locked out of your Google account, what do you do? (2021)",
"position": 41
},
{
"author": "db48x",
"karma": 3876,
"comments": 369,
"title": "The best engineering interview question I've ever gotten",
"position": 8
},
{
"author": "bpierre",
"karma": 42339,
"comments": 98,
"title": "No one knows why the most used spacecraft propulsion system works",
"position": 39
},
{
"author": "wglb",
"karma": 50533,
"comments": 115,
"title": "James Webb telescope's coldest instrument reaches operating temperature",
"position": 32
},
{
"author": "elsewhen",
"karma": 19297,
"comments": 382,
"title": "Americans are drowning in spam",
"position": 7
},
{
"author": "_wldu",
"karma": 5058,
"comments": 177,
"title": "Assume your devices are compromised",
"position": 19
},
{
"author": "igonvalue",
"karma": 2689,
"comments": 131,
"title": "The casualties at the other end of the remote-controlled kill",
"position": 31
},
{
"author": "ingve",
"karma": 155403,
"comments": 247,
"title": "M1 Thunderbolt ports don’t fully support USB 3.1 Gen 2",
"position": 13
},
{
"author": "lkrubner",
"karma": 13368,
"comments": 139,
"title": "Shirky.com is gone",
"position": 25
},
{
"author": "luu",
"karma": 86882,
"comments": 98,
"title": "The games Nintendo didn't want you to play: Tengen",
"position": 38
},
{
"author": "pseudolus",
"karma": 118368,
"comments": 224,
"title": "Jack Dorsey’s $2.9M NFT dropped 99% in value",
"position": 16
},
{
"author": "MilnerRoute",
"karma": 11952,
"comments": 80,
"title": "California teen with autism who vanished 3 years ago is found alive in Utah",
"position": 47
},
{
"author": "maxerickson",
"karma": 33323,
"comments": 508,
"title": "The Uber Bubble",
"position": 4
},
{
"author": "zeristor",
"karma": 6647,
"comments": 174,
"title": "Barbary Pirates and English Slaves (2017)",
"position": 20
},
{
"author": "walterbell",
"karma": 60391,
"comments": 86,
"title": "Eric Schmidt's influence on U.S. science policy",
"position": 45
},
{
"author": "diego",
"karma": 8235,
"comments": 287,
"title": "The money I saved as a child would buy one picogram of gold today",
"position": 9
},
{
"author": "AndrewDucker",
"karma": 30341,
"comments": 225,
"title": "Notable items missing from English Wikipedia",
"position": 15
},
{
"author": "mgh2",
"karma": 5052,
"comments": 78,
"title": "How Rainy Is Seattle? It's Not Even in the Top of Major U.S. Cities (2019)",
"position": 48
},
{
"author": "optimalsolver",
"karma": 5373,
"comments": 134,
"title": "List Of Adhesive Tapes",
"position": 30
},
{
"author": "grawprog",
"karma": 10654,
"comments": 102,
"title": "Vancouver proposes 10k-dollar annual fee for gas stations without EV charging",
"position": 36
},
{
"author": "mzs",
"karma": 6691,
"comments": 168,
"title": "Ten members of international stock manipulation ring charged in Manhattan",
"position": 22
},
{
"author": "Anon84",
"karma": 42095,
"comments": 139,
"title": "The Principles of Deep Learning Theory",
"position": 26
},
{
"author": "tosh",
"karma": 109351,
"comments": 87,
"title": "27 years ago I accidentally ran the hardest, strangest Easter egg hunt",
"position": 44
},
{
"author": "rustoo",
"karma": 10558,
"comments": 173,
"title": "A drug that cures alcoholism may be the next anti-anxiety medication",
"position": 21
},
{
"author": "WoodenChair",
"karma": 5932,
"comments": 531,
"title": "Richard Stallman – The state of the Free Software movement",
"position": 3
},
{
"author": "kaycebasques",
"karma": 4696,
"comments": 87,
"title": "The chicken you are eating has increased 364% in size over the last 50 years",
"position": 43
},
{
"author": "slg",
"karma": 25678,
"comments": 615,
"title": "My take on Elon's offer for Twitter",
"position": 1
},
{
"author": "germinalphrase",
"karma": 4284,
"comments": 89,
"title": "Primer: Statistical Armour",
"position": 42
},
{
"author": "exolymph",
"karma": 11235,
"comments": 151,
"title": "Collectibles are terrible investments",
"position": 24
},
{
"author": "RickJWagner",
"karma": 7297,
"comments": 135,
"title": "100 People with rare cancers who attended same NJ high school demand answers",
"position": 28
},
{
"author": "mooreds",
"karma": 38643,
"comments": 607,
"title": "The Colorado Safety Stop is the law of the land",
"position": 2
},
{
"author": "TangerineDream",
"karma": 8198,
"comments": 431,
"title": "DuckDuckGo Removes Pirate Sites and YouTube-DL from Its Search Results",
"position": 5
},
{
"author": "maxerickson",
"karma": 33323,
"comments": 210,
"title": "Reversing hearing loss with regenerative therapy",
"position": 17
},
{
"author": "DantesKite",
"karma": 2416,
"comments": 155,
"title": "Psychedelics and mental illness",
"position": 23
},
{
"author": "lukastyrychtr",
"karma": 4943,
"comments": 93,
"title": "Pointers Are Complicated III, or: Pointer-integer casts exposed",
"position": 40
},
{
"author": "tosh",
"karma": 109351,
"comments": 85,
"title": "Seashore: Easy to use Mac OS X image editing application for the rest of us",
"position": 46
},
{
"author": "kiyanwang",
"karma": 22866,
"comments": 230,
"title": "I Avoid Async/Await",
"position": 14
},
{
"author": "pcr910303",
"karma": 22957,
"comments": 262,
"title": "I hope distributed is not the new default",
"position": 12
},
{
"author": "ilamont",
"karma": 40680,
"comments": 273,
"title": "The silenced deaths of the Shanghai 2022 lockdown",
"position": 11
},
{
"author": "rcarmo",
"karma": 13218,
"comments": 422,
"title": "It’s Still Stupidly, Difficult to Buy a ‘Dumb’ TV",
"position": 6
},
{
"author": "dshipper",
"karma": 4354,
"comments": 135,
"title": "Twitter Should Open Up the Algorithm",
"position": 29
},
{
"author": "akanet",
"karma": 2421,
"comments": 111,
"title": "Single mom sues coding boot camp over job placement rates",
"position": 33
},
{
"author": "Tomte",
"karma": 101902,
"comments": 103,
"title": "Birds Aren’t Real took on the conspiracy theorists",
"position": 35
},
{
"author": "mikestew",
"karma": 18007,
"comments": 206,
"title": "How destructive are nuclear weapons really?",
"position": 18
},
{
"author": "prostoalex",
"karma": 114995,
"comments": 78,
"title": "Two Silicon Valley executives charged with H-1B visa fraud",
"position": 49
},
{
"author": "prostoalex",
"karma": 114995,
"comments": 136,
"title": "As climate fears mount, some are relocating within the US",
"position": 27
},
{
"author": "nixass",
"karma": 5294,
"comments": 102,
"title": "Toyota warns about rushing into electrification",
"position": 37
},
{
"author": "DyslexicAtheist",
"karma": 32296,
"comments": 110,
"title": "Study on US-Russia nuclear war: 91.5 MM casualties in first few hours (2019)",
"position": 34
}
]
```

2. `GET /health`

```json
{
"status": "UP"
}
```

3. `GET /metrics`

```
# HELP TotalTime Total request latency time taken by each request
# TYPE TotalTime gauge
TotalTime 0
# HELP go_gc_cycles_automatic_gc_cycles_total Count of completed GC cycles generated by the Go runtime.
# TYPE go_gc_cycles_automatic_gc_cycles_total counter
go_gc_cycles_automatic_gc_cycles_total 25
# HELP go_gc_cycles_forced_gc_cycles_total Count of completed GC cycles forced by the application.
# TYPE go_gc_cycles_forced_gc_cycles_total counter
go_gc_cycles_forced_gc_cycles_total 0
# HELP go_gc_cycles_total_gc_cycles_total Count of all completed GC cycles.
# TYPE go_gc_cycles_total_gc_cycles_total counter
go_gc_cycles_total_gc_cycles_total 25
# HELP go_gc_duration_seconds A summary of the pause duration of garbage collection cycles.
# TYPE go_gc_duration_seconds summary
go_gc_duration_seconds{quantile="0"} 6.81e-05
go_gc_duration_seconds{quantile="0.25"} 0.0002321
go_gc_duration_seconds{quantile="0.5"} 0.0003511
go_gc_duration_seconds{quantile="0.75"} 0.0008538
go_gc_duration_seconds{quantile="1"} 0.0025279
go_gc_duration_seconds_sum 0.0157166
go_gc_duration_seconds_count 25
# HELP go_gc_heap_allocs_by_size_bytes_total Distribution of heap allocations by approximate size. Note that this does not include tiny objects as defined by /gc/heap/tiny/allocs:objects, only tiny blocks.
# TYPE go_gc_heap_allocs_by_size_bytes_total histogram
go_gc_heap_allocs_by_size_bytes_total_bucket{le="8.999999999999998"} 8446
go_gc_heap_allocs_by_size_bytes_total_bucket{le="24.999999999999996"} 901146
go_gc_heap_allocs_by_size_bytes_total_bucket{le="64.99999999999999"} 1.046863e+06
go_gc_heap_allocs_by_size_bytes_total_bucket{le="144.99999999999997"} 1.1366e+06
go_gc_heap_allocs_by_size_bytes_total_bucket{le="320.99999999999994"} 1.160712e+06
go_gc_heap_allocs_by_size_bytes_total_bucket{le="704.9999999999999"} 1.179908e+06
go_gc_heap_allocs_by_size_bytes_total_bucket{le="1536.9999999999998"} 1.190926e+06
go_gc_heap_allocs_by_size_bytes_total_bucket{le="3200.9999999999995"} 1.194572e+06
go_gc_heap_allocs_by_size_bytes_total_bucket{le="6528.999999999999"} 1.197874e+06
go_gc_heap_allocs_by_size_bytes_total_bucket{le="13568.999999999998"} 1.198659e+06
go_gc_heap_allocs_by_size_bytes_total_bucket{le="27264.999999999996"} 1.199153e+06
go_gc_heap_allocs_by_size_bytes_total_bucket{le="+Inf"} 1.200027e+06
go_gc_heap_allocs_by_size_bytes_total_sum 1.62226088e+08
go_gc_heap_allocs_by_size_bytes_total_count 1.200027e+06
# HELP go_gc_heap_allocs_bytes_total Cumulative sum of memory allocated to the heap by the application.
# TYPE go_gc_heap_allocs_bytes_total counter
go_gc_heap_allocs_bytes_total 1.62226088e+08
# HELP go_gc_heap_allocs_objects_total Cumulative count of heap allocations triggered by the application. Note that this does not include tiny objects as defined by /gc/heap/tiny/allocs:objects, only tiny blocks.
# TYPE go_gc_heap_allocs_objects_total counter
go_gc_heap_allocs_objects_total 1.200027e+06
# HELP go_gc_heap_frees_by_size_bytes_total Distribution of freed heap allocations by approximate size. Note that this does not include tiny objects as defined by /gc/heap/tiny/allocs:objects, only tiny blocks.
# TYPE go_gc_heap_frees_by_size_bytes_total histogram
go_gc_heap_frees_by_size_bytes_total_bucket{le="8.999999999999998"} 6318
go_gc_heap_frees_by_size_bytes_total_bucket{le="24.999999999999996"} 838564
go_gc_heap_frees_by_size_bytes_total_bucket{le="64.99999999999999"} 972917
go_gc_heap_frees_by_size_bytes_total_bucket{le="144.99999999999997"} 1.05744e+06
go_gc_heap_frees_by_size_bytes_total_bucket{le="320.99999999999994"} 1.079728e+06
go_gc_heap_frees_by_size_bytes_total_bucket{le="704.9999999999999"} 1.097584e+06
go_gc_heap_frees_by_size_bytes_total_bucket{le="1536.9999999999998"} 1.108002e+06
go_gc_heap_frees_by_size_bytes_total_bucket{le="3200.9999999999995"} 1.11151e+06
go_gc_heap_frees_by_size_bytes_total_bucket{le="6528.999999999999"} 1.114463e+06
go_gc_heap_frees_by_size_bytes_total_bucket{le="13568.999999999998"} 1.115187e+06
go_gc_heap_frees_by_size_bytes_total_bucket{le="27264.999999999996"} 1.115659e+06
go_gc_heap_frees_by_size_bytes_total_bucket{le="+Inf"} 1.116452e+06
go_gc_heap_frees_by_size_bytes_total_sum 1.4936852e+08
go_gc_heap_frees_by_size_bytes_total_count 1.116452e+06
# HELP go_gc_heap_frees_bytes_total Cumulative sum of heap memory freed by the garbage collector.
# TYPE go_gc_heap_frees_bytes_total counter
go_gc_heap_frees_bytes_total 1.4936852e+08
# HELP go_gc_heap_frees_objects_total Cumulative count of heap allocations whose storage was freed by the garbage collector. Note that this does not include tiny objects as defined by /gc/heap/tiny/allocs:objects, only tiny blocks.
# TYPE go_gc_heap_frees_objects_total counter
go_gc_heap_frees_objects_total 1.116452e+06
# HELP go_gc_heap_goal_bytes Heap size target for the end of the GC cycle.
# TYPE go_gc_heap_goal_bytes gauge
go_gc_heap_goal_bytes 1.3706784e+07
# HELP go_gc_heap_objects_objects Number of objects, live or unswept, occupying heap memory.
# TYPE go_gc_heap_objects_objects gauge
go_gc_heap_objects_objects 83575
# HELP go_gc_heap_tiny_allocs_objects_total Count of small allocations that are packed together into blocks. These allocations are counted separately from other allocations because each individual allocation is not tracked by the runtime, only their block. Each block is already accounted for in allocs-by-size and frees-by-size.
# TYPE go_gc_heap_tiny_allocs_objects_total counter
go_gc_heap_tiny_allocs_objects_total 760226
# HELP go_gc_pauses_seconds_total Distribution individual GC-related stop-the-world pause latencies.
# TYPE go_gc_pauses_seconds_total histogram
go_gc_pauses_seconds_total_bucket{le="-5e-324"} 0
go_gc_pauses_seconds_total_bucket{le="9.999999999999999e-10"} 0
go_gc_pauses_seconds_total_bucket{le="9.999999999999999e-09"} 0
go_gc_pauses_seconds_total_bucket{le="9.999999999999998e-08"} 0
go_gc_pauses_seconds_total_bucket{le="1.0239999999999999e-06"} 0
go_gc_pauses_seconds_total_bucket{le="1.0239999999999999e-05"} 0
go_gc_pauses_seconds_total_bucket{le="0.00010239999999999998"} 18
go_gc_pauses_seconds_total_bucket{le="0.0010485759999999998"} 47
go_gc_pauses_seconds_total_bucket{le="0.010485759999999998"} 50
go_gc_pauses_seconds_total_bucket{le="0.10485759999999998"} 50
go_gc_pauses_seconds_total_bucket{le="+Inf"} 50
go_gc_pauses_seconds_total_sum NaN
go_gc_pauses_seconds_total_count 50
# HELP go_goroutines Number of goroutines that currently exist.
# TYPE go_goroutines gauge
go_goroutines 8
# HELP go_info Information about the Go environment.
# TYPE go_info gauge
go_info{version="go1.18.1"} 1
# HELP go_memory_classes_heap_free_bytes Memory that is completely free and eligible to be returned to the underlying system, but has not been. This metric is the runtime's estimate of free address space that is backed by physical memory.
# TYPE go_memory_classes_heap_free_bytes gauge
go_memory_classes_heap_free_bytes 761856
# HELP go_memory_classes_heap_objects_bytes Memory occupied by live objects and dead objects that have not yet been marked free by the garbage collector.
# TYPE go_memory_classes_heap_objects_bytes gauge
go_memory_classes_heap_objects_bytes 1.2857568e+07
# HELP go_memory_classes_heap_released_bytes Memory that is completely free and has been returned to the underlying system. This metric is the runtime's estimate of free address space that is still mapped into the process, but is not backed by physical memory.
# TYPE go_memory_classes_heap_released_bytes gauge
go_memory_classes_heap_released_bytes 1.0346496e+07
# HELP go_memory_classes_heap_stacks_bytes Memory allocated from the heap that is reserved for stack space, whether or not it is currently in-use.
# TYPE go_memory_classes_heap_stacks_bytes gauge
go_memory_classes_heap_stacks_bytes 1.376256e+06
# HELP go_memory_classes_heap_unused_bytes Memory that is reserved for heap objects but is not currently used to hold heap objects.
# TYPE go_memory_classes_heap_unused_bytes gauge
go_memory_classes_heap_unused_bytes 4.017952e+06
# HELP go_memory_classes_metadata_mcache_free_bytes Memory that is reserved for runtime mcache structures, but not in-use.
# TYPE go_memory_classes_metadata_mcache_free_bytes gauge
go_memory_classes_metadata_mcache_free_bytes 10800
# HELP go_memory_classes_metadata_mcache_inuse_bytes Memory that is occupied by runtime mcache structures that are currently being used.
# TYPE go_memory_classes_metadata_mcache_inuse_bytes gauge
go_memory_classes_metadata_mcache_inuse_bytes 4800
# HELP go_memory_classes_metadata_mspan_free_bytes Memory that is reserved for runtime mspan structures, but not in-use.
# TYPE go_memory_classes_metadata_mspan_free_bytes gauge
go_memory_classes_metadata_mspan_free_bytes 45016
# HELP go_memory_classes_metadata_mspan_inuse_bytes Memory that is occupied by runtime mspan structures that are currently being used.
# TYPE go_memory_classes_metadata_mspan_inuse_bytes gauge
go_memory_classes_metadata_mspan_inuse_bytes 199784
# HELP go_memory_classes_metadata_other_bytes Memory that is reserved for or used to hold runtime metadata.
# TYPE go_memory_classes_metadata_other_bytes gauge
go_memory_classes_metadata_other_bytes 5.64524e+06
# HELP go_memory_classes_os_stacks_bytes Stack memory allocated by the underlying operating system.
# TYPE go_memory_classes_os_stacks_bytes gauge
go_memory_classes_os_stacks_bytes 0
# HELP go_memory_classes_other_bytes Memory used by execution trace buffers, structures for debugging the runtime, finalizer and profiler specials, and more.
# TYPE go_memory_classes_other_bytes gauge
go_memory_classes_other_bytes 859061
# HELP go_memory_classes_profiling_buckets_bytes Memory that is used by the stack trace hash map used for profiling.
# TYPE go_memory_classes_profiling_buckets_bytes gauge
go_memory_classes_profiling_buckets_bytes 4723
# HELP go_memory_classes_total_bytes All memory mapped by the Go runtime into the current process as read-write. Note that this does not include memory mapped by code called via cgo or via the syscall package. Sum of all metrics in /memory/classes.
# TYPE go_memory_classes_total_bytes gauge
go_memory_classes_total_bytes 3.6129552e+07
# HELP go_memstats_alloc_bytes Number of bytes allocated and still in use.
# TYPE go_memstats_alloc_bytes gauge
go_memstats_alloc_bytes 1.2857568e+07
# HELP go_memstats_alloc_bytes_total Total number of bytes allocated, even if freed.
# TYPE go_memstats_alloc_bytes_total counter
go_memstats_alloc_bytes_total 1.62226088e+08
# HELP go_memstats_buck_hash_sys_bytes Number of bytes used by the profiling bucket hash table.
# TYPE go_memstats_buck_hash_sys_bytes gauge
go_memstats_buck_hash_sys_bytes 4723
# HELP go_memstats_frees_total Total number of frees.
# TYPE go_memstats_frees_total counter
go_memstats_frees_total 1.876678e+06
# HELP go_memstats_gc_cpu_fraction The fraction of this program's available CPU time used by the GC since the program started.
# TYPE go_memstats_gc_cpu_fraction gauge
go_memstats_gc_cpu_fraction 0
# HELP go_memstats_gc_sys_bytes Number of bytes used for garbage collection system metadata.
# TYPE go_memstats_gc_sys_bytes gauge
go_memstats_gc_sys_bytes 5.64524e+06
# HELP go_memstats_heap_alloc_bytes Number of heap bytes allocated and still in use.
# TYPE go_memstats_heap_alloc_bytes gauge
go_memstats_heap_alloc_bytes 1.2857568e+07
# HELP go_memstats_heap_idle_bytes Number of heap bytes waiting to be used.
# TYPE go_memstats_heap_idle_bytes gauge
go_memstats_heap_idle_bytes 1.1108352e+07
# HELP go_memstats_heap_inuse_bytes Number of heap bytes that are in use.
# TYPE go_memstats_heap_inuse_bytes gauge
go_memstats_heap_inuse_bytes 1.687552e+07
# HELP go_memstats_heap_objects Number of allocated objects.
# TYPE go_memstats_heap_objects gauge
go_memstats_heap_objects 83575
# HELP go_memstats_heap_released_bytes Number of heap bytes released to OS.
# TYPE go_memstats_heap_released_bytes gauge
go_memstats_heap_released_bytes 1.0346496e+07
# HELP go_memstats_heap_sys_bytes Number of heap bytes obtained from system.
# TYPE go_memstats_heap_sys_bytes gauge
go_memstats_heap_sys_bytes 2.7983872e+07
# HELP go_memstats_last_gc_time_seconds Number of seconds since 1970 of last garbage collection.
# TYPE go_memstats_last_gc_time_seconds gauge
go_memstats_last_gc_time_seconds 1.6503106247604249e+09
# HELP go_memstats_lookups_total Total number of pointer lookups.
# TYPE go_memstats_lookups_total counter
go_memstats_lookups_total 0
# HELP go_memstats_mallocs_total Total number of mallocs.
# TYPE go_memstats_mallocs_total counter
go_memstats_mallocs_total 1.960253e+06
# HELP go_memstats_mcache_inuse_bytes Number of bytes in use by mcache structures.
# TYPE go_memstats_mcache_inuse_bytes gauge
go_memstats_mcache_inuse_bytes 4800
# HELP go_memstats_mcache_sys_bytes Number of bytes used for mcache structures obtained from system.
# TYPE go_memstats_mcache_sys_bytes gauge
go_memstats_mcache_sys_bytes 15600
# HELP go_memstats_mspan_inuse_bytes Number of bytes in use by mspan structures.
# TYPE go_memstats_mspan_inuse_bytes gauge
go_memstats_mspan_inuse_bytes 199784
# HELP go_memstats_mspan_sys_bytes Number of bytes used for mspan structures obtained from system.
# TYPE go_memstats_mspan_sys_bytes gauge
go_memstats_mspan_sys_bytes 244800
# HELP go_memstats_next_gc_bytes Number of heap bytes when next garbage collection will take place.
# TYPE go_memstats_next_gc_bytes gauge
go_memstats_next_gc_bytes 1.3706784e+07
# HELP go_memstats_other_sys_bytes Number of bytes used for other system allocations.
# TYPE go_memstats_other_sys_bytes gauge
go_memstats_other_sys_bytes 859061
# HELP go_memstats_stack_inuse_bytes Number of bytes in use by the stack allocator.
# TYPE go_memstats_stack_inuse_bytes gauge
go_memstats_stack_inuse_bytes 1.376256e+06
# HELP go_memstats_stack_sys_bytes Number of bytes obtained from system for stack allocator.
# TYPE go_memstats_stack_sys_bytes gauge
go_memstats_stack_sys_bytes 1.376256e+06
# HELP go_memstats_sys_bytes Number of bytes obtained from system.
# TYPE go_memstats_sys_bytes gauge
go_memstats_sys_bytes 3.6129552e+07
# HELP go_sched_goroutines_goroutines Count of live goroutines.
# TYPE go_sched_goroutines_goroutines gauge
go_sched_goroutines_goroutines 7
# HELP go_sched_latencies_seconds Distribution of the time goroutines have spent in the scheduler in a runnable state before actually running.
# TYPE go_sched_latencies_seconds histogram
go_sched_latencies_seconds_bucket{le="-5e-324"} 0
go_sched_latencies_seconds_bucket{le="9.999999999999999e-10"} 2985
go_sched_latencies_seconds_bucket{le="9.999999999999999e-09"} 2985
go_sched_latencies_seconds_bucket{le="9.999999999999998e-08"} 2985
go_sched_latencies_seconds_bucket{le="1.0239999999999999e-06"} 2985
go_sched_latencies_seconds_bucket{le="1.0239999999999999e-05"} 3742
go_sched_latencies_seconds_bucket{le="0.00010239999999999998"} 7246
go_sched_latencies_seconds_bucket{le="0.0010485759999999998"} 7905
go_sched_latencies_seconds_bucket{le="0.010485759999999998"} 8072
go_sched_latencies_seconds_bucket{le="0.10485759999999998"} 8079
go_sched_latencies_seconds_bucket{le="+Inf"} 8079
go_sched_latencies_seconds_sum NaN
go_sched_latencies_seconds_count 8079
# HELP go_threads Number of OS threads created.
# TYPE go_threads gauge
go_threads 9
# HELP process_cpu_seconds_total Total user and system CPU time spent in seconds.
# TYPE process_cpu_seconds_total counter
process_cpu_seconds_total 4.29
# HELP process_max_fds Maximum number of open file descriptors.
# TYPE process_max_fds gauge
process_max_fds 1.048576e+06
# HELP process_open_fds Number of open file descriptors.
# TYPE process_open_fds gauge
process_open_fds 9
# HELP process_resident_memory_bytes Resident memory size in bytes.
# TYPE process_resident_memory_bytes gauge
process_resident_memory_bytes 3.217408e+07
# HELP process_start_time_seconds Start time of the process since unix epoch in seconds.
# TYPE process_start_time_seconds gauge
process_start_time_seconds 1.65031061401e+09
# HELP process_virtual_memory_bytes Virtual memory size in bytes.
# TYPE process_virtual_memory_bytes gauge
process_virtual_memory_bytes 7.33421568e+08
# HELP process_virtual_memory_max_bytes Maximum amount of virtual memory available in bytes.
# TYPE process_virtual_memory_max_bytes gauge
process_virtual_memory_max_bytes 1.8446744073709552e+19
# HELP promhttp_metric_handler_requests_in_flight Current number of scrapes being served.
# TYPE promhttp_metric_handler_requests_in_flight gauge
promhttp_metric_handler_requests_in_flight 1
# HELP promhttp_metric_handler_requests_total Total number of scrapes by HTTP status code.
# TYPE promhttp_metric_handler_requests_total counter
promhttp_metric_handler_requests_total{code="200"} 0
promhttp_metric_handler_requests_total{code="500"} 0
promhttp_metric_handler_requests_total{code="503"} 0
```

# Response time and request count

In the code we are explicitely exposing the latency time and the number of requests are tracking via prometheus
```
[GIN-debug] Environment variable PORT is undefined. Using port :8080 by default
[GIN-debug] Listening and serving HTTP on :8080
[GIN] 2022/04/18 - 19:36:58 | 404 |        25.9µs |      172.17.0.1 | GET      "/"
7:37PM INF  TTFB=414.8473 TotalTime=414.9877 URL=https://hacker-news.firebaseio.com/v0/topstories.json
7:37PM INF  TTFB=111.6459 TotalTime=111.8826 URL=https://hacker-news.firebaseio.com/v0/item/31072590.json
7:37PM INF  TTFB=224.45 TotalTime=225.1808 URL=https://hacker-news.firebaseio.com/v0/item/31074177.json
7:37PM INF  TTFB=340.041 TotalTime=340.1302 URL=https://hacker-news.firebaseio.com/v0/item/31072485.json
7:37PM INF  TTFB=345.39 TotalTime=345.4917 URL=https://hacker-news.firebaseio.com/v0/item/31068479.json
7:37PM INF  TTFB=359.4824 TotalTime=359.5826 URL=https://hacker-news.firebaseio.com/v0/item/31074896.json
```

Total request counts

```
7:37PM INF sent /stories TotalTime=3764.1212
[GIN] 2022/04/18 - 19:37:04 | 200 |    3.7643643s |      172.17.0.1 | GET      "/stories"
[GIN] 2022/04/18 - 19:37:46 | 200 |        78.1µs |      172.17.0.1 | GET      "/health"
[GIN] 2022/04/18 - 19:38:13 | 200 |      6.7153ms |      172.17.0.1 | GET      "/metrics"
```

# Further Improvements

1. Have added unit tests for some services. Need to add more
2. Have enabled prometheus for this app and it provides basic golang metrics. Need to extend the latency based performance metrics. As of now the logs have the request latency durations
3. Currently, I am using 50 workers for sorting the 500 stories. Not sure about the rate limiting values for hacker news api. But we can definitely add 50 more workers and improve the query performance
4. Improve caching and use an in memory cache like redis/mongodb for faster execution at scale, lets say we need to make this a production grade appliocation. As of now I am using the golang cache for faster execution
5. I am exposing the service as a node port and it does not leverage externalDNS and certs. The helm charts have those dependencies but as of now those are disabled

