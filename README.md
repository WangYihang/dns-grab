# DNS Grab

## Description

`dns-grab` is a tool for doing DNS queries from a list of domain names.

## Installation

```bash
go install github.com/WangYihang/dns-grab@latest
```

## Usage

```bash
$ dns-grab --h
Usage:
  dns-grab [OPTIONS]

Application Options:
      --input=                        input file path
      --output=                       output file path
      --status=                       status file path (default: -)
      --num-workers=                  number of workers (default: 32)
      --num-shards=                   number of shards (default: 1)
      --shard=                        shard (default: 0)
      --max-tries=                    max tries (default: 2)
      --max-runtime-per-task-seconds= max runtime per task seconds (default: 8)
      --qtype=                        qtype (default: A)
      --resolver=                     resolver (default: 8.8.8.8:53)
      --version                       print version and exit

Help Options:
  -h, --help                          Show this help message

```

```bash
$ head input.txt
www.baidu.com
www.jd.com
www.tsinghua.edu.cn
www.pku.edu.cn
```

```bash
$ dns-grab --input input.txt --output output.txt
...
```

```bash
$ head -n 1 output.txt | jq
```

```json
{
  "index": 0,
  "id": "3a46f868-57e8-4ebe-bd4f-5d983e9f6a4c",
  "started_at": 1712683765299936,
  "finished_at": 1712683765300326,
  "num_tries": 1,
  "task": {
    "qname": "www.baidu.com",
    "qtype": "A",
    "resolver": "8.8.8.8:53",
    "dns": {
      "request": {
        "header": {
          "id": 10740,
          "qr": false,
          "op": "QUERY",
          "aa": false,
          "tc": false,
          "rd": true,
          "ra": false,
          "z": false,
          "ad": false,
          "cd": false,
          "rc": "NOERROR"
        },
        "questions": [
          {
            "name": "www.baidu.com.",
            "qtype": "A",
            "qclass": "IN"
          }
        ]
      },
      "response": {
        "header": {
          "id": 10740,
          "qr": true,
          "op": "QUERY",
          "aa": false,
          "tc": false,
          "rd": true,
          "ra": true,
          "z": false,
          "ad": false,
          "cd": false,
          "rc": "NOERROR"
        },
        "questions": [
          {
            "name": "www.baidu.com.",
            "qtype": "A",
            "qclass": "IN"
          }
        ],
        "answers": [
          {
            "rtype": "CNAME",
            "rname": "www.baidu.com.",
            "cname": "www.a.shifen.com."
          },
          {
            "rtype": "A",
            "rname": "www.a.shifen.com.",
            "a": "182.61.200.6"
          },
          {
            "rtype": "A",
            "rname": "www.a.shifen.com.",
            "a": "182.61.200.7"
          }
        ]
      }
    }
  },
  "error": ""
}
```