from diagrams import Cluster, Diagram
from diagrams.aws.compute import ECS
from diagrams.aws.database import ElastiCache, DB
from diagrams.aws.network import ELB
from diagrams.aws.network import APIGateway

with Diagram("LinkZapURL Architecture", show=False):

    with Cluster("Nginx Service"):
        gateway = APIGateway("API Gateway")
        lb = ELB("Load Balancer")

    with Cluster("Go Service"):
        svc_group = [
            ECS("server-1"),
            ECS("server-2"),
        ]

    with Cluster("MongoDB Cluster"):
        mongo_router_group = [
            DB("router-1"),
            DB("router-2"),
        ]

        with Cluster("MongoDB Shard"):
            mongo_shard_group = [
                DB("shard-1"),
                DB("shard-2"),
                DB("shard-3"),
            ]

        mongo_router_group[0] >> mongo_shard_group[0]

    with Cluster("Redis Cluster"):
        redis_nodes = [
            ElastiCache(f"node {i}") for i in range(1, 7)
            ]

    gateway >> lb >> svc_group
    svc_group >> mongo_router_group[0]

    svc_group >> redis_nodes[0]
    
