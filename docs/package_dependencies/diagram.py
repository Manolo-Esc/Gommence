
# Available logos: https://diagrams.mingrammer.com/docs/nodes/onprem

from diagrams import Diagram, Cluster, Edge, Node
#from diagrams.c4 import Container
# from diagrams.generic.network import Router, Switch
# # diagrams.generic.storage.Storage
# from diagrams.aws.media import MediaServices
# from diagrams.aws.network import RouteTable
# from diagrams.custom import Custom
# from diagrams.onprem.container import Docker, K3S
# from diagrams.onprem.database import Postgresql
# from diagrams.onprem.monitoring import Grafana, Prometheus
# from diagrams.onprem.network import Nginx
# from diagrams.onprem.proxmox import Pve
# from diagrams.onprem.vcs import Git
# from diagrams.programming.language import Swift, Php

#from diagrams.programming.language import Go
from diagrams.azure.compute import VMLinux as Go

#help(Diagram)
#help(Cluster)
#help(Docker)
#help(Edge)

graph_attr = {"layout": "dot"      # default: dot
              # No parecen tener efecto:
              #"nodesep": "0.05" 
              }   
smallNode = {"fixedsize": "true", "width": "0.6", "height": "0.6"}
#tinyNode  = {"fixedsize": "true", "width": "0.3", "height": "0.3"}
bigNode   = {"fixedsize": "true", "width": "2", "height": "3.8"} 

smallEdge = {
             "minlen": "1", # no funciona bien 
             "len": "0.001"  # not used in layout engine `dot``
             # No parecen tener efecto:
             #"showboxes": "true",
             #"dir": "none", 
             } 
verticalizeEdge = {"minlen": "0.5"} # parece poner en vertical la orientacion de los nodos que une

ports_forwardings = (
   "80   → 192.168.10.11\n"
   "443  → 192.168.10.11\n"
   "5080 → 192.168.10.13\n"
   "1935 → 192.168.10.13\n"
   "4200 udp → 192..10.13\n"
   "50000-6000 udp → ...13"
)

def InvEdge():
    return Edge(style="invis")


#with Diagram("", filename="silentbob", show=True, direction="LR", node_attr={"fixedsize": "true", "width": "1", "height": "1"}):
with Diagram("", filename="diagram", show=True, direction="LR", graph_attr=graph_attr):
    cmd = Go("main")
    server = Go("server")
    repos = Go("repos")
    rest = Go("rest")
    app = Go("app")
    domain = Go("domain")
    dtos = Go("dtos")
    database = Go("db migration")
    jwt = Go("jwt")
    uid = Go("uid")
    mocks = Go("mocks")
    ports = Go("ports")
    cache = Go("cache")
    logger = Go("logger")
    network = Go("network")
    otel = Go("open telemetry")
    validator = Go("validator")

    # XXX tests end2end, integration

    cmd >> server
    server >> rest
    server >> logger
    server >> cache
    server >> otel
    server >> app
    server >> ports
    server >> repos
    ports >> dtos
    ports >> domain
    dtos >> domain
    app >> ports
    app >> dtos
    app >> domain
    rest >> ports
    rest >> dtos
    repos >> dtos
    repos >> ports
    repos >> domain


    # telefonica = Router("Telefonica\n 213.97.87.110\n Modo Bridge ") 
    # Vodafone = Router("Vodafone\n IP variable\n 10.0.1.1") 

    # with Cluster("Red oficina"):
    #     oficina = Switch("Switch", **{"fixedsize": "true", "width": "2"})

    # with Cluster("Unify Gateway"):
    #     unify = Switch("Unify\n192.168.10.1", **{"fixedsize": "true", "width": "2", "height": "1.5"})
    #     telefonica - unify
    #     telefonica >> Edge(style="invis", **verticalizeEdge) >> RouteTable(ports_forwardings) #, **{"fixedsize": "true", "height": "5.5"})

    # with Cluster("SilentBob Bare metal"):
    #         # silentbob =  Pve("Proxmox\n 10.0.1.6:8006 \n root \n dual1234")
    #         # silentbob >> Edge(style="invis", **trickEdge) >> Custom("silentBob\n 10.0.1.177 \n ADMIN \n PAVQEJMQHX \n vmbr1: 192.168.10.10\n vmbr0: 10.0.1.6", "./resources/server.png", **bigNode)
    #         silentbob_pc =  Custom("silentBob\n 10.0.1.177 \n ADMIN \n PAVQEJMQHX \n vmbr1: 192.168.10.10\n vmbr0: 10.0.1.6", "./resources/server.png", **bigNode)
    #         proxmox = Pve("Proxmox\n 10.0.1.6:8006 \n root \n dual1234")
    #         silentbob = proxmox
    #         silentbob_pc >> Edge(style="invis", **verticalizeEdge) >> proxmox

    # with Cluster("SilentBob VMs"):
        

    # #with Cluster("", graph_attr={"rank": "same"}):
    #     with Cluster("v3"):
    #         v3 = VMLinux("v3\n 192.168.10.11\n 10.0.1.7 \n dualuser \n 1234")
    #         silentbob >> InvEdge() >> v3

    #         with Cluster("Observability"):
    #             with Cluster("Node exporter"):
    #                 Docker("", **smallNode) >> InvEdge() >> Custom("", "./resources/prom-exp.png", **smallNode)
    #             with Cluster("Blackbox exporter (services & URLs)"):
    #                 v3 >> InvEdge() >> Docker("", **smallNode) >> InvEdge() >> Custom("", "./resources/prom-exp.png", **smallNode)
    #             with Cluster("cAdvisor exporter (containers)"):
    #                 #v3 >> InvEdge() >> Docker("", **smallNode) >> Edge(style="invis", **smallEdge) >> Custom("", "./resources/prom-exp.png", **smallNode)
    #                 v3 >> InvEdge() >> Docker("", **smallNode) >> Edge(style="invis") >> Custom("", "./resources/prom-exp.png", **smallNode)
    #             with Cluster("Nginx exporter"):
    #                 #v3 >> InvEdge() >> Docker("", **smallNode) >> Edge(style="invis", **smallEdge) >> Custom("", "./resources/prom-exp.png", **smallNode)
    #                 v3 >> InvEdge() >> Docker("", **smallNode) >> Edge(style="invis") >> Custom("", "./resources/prom-exp.png", **smallNode)
    #             with Cluster("Postgres exporter"):
    #                 #v3 >> InvEdge() >> Docker("", **smallNode) >> Edge(style="invis", **smallEdge) >> Custom("", "./resources/prom-exp.png", **smallNode)
    #                 v3 >> InvEdge() >> Docker("", **smallNode) >> Edge(style="invis") >> Custom("", "./resources/prom-exp.png", **smallNode)
    #             with Cluster("prometheus: 10.0.1.7:9090"):
    #                 #v3 >> InvEdge() >> Docker("", **smallNode) >> Edge(style="invis", **smallEdge) >> Prometheus("", **smallNode)
    #                 v3 >> InvEdge() >> Docker("", **smallNode) >> Edge(style="invis") >> Prometheus("", **smallNode)
    #             with Cluster("grafana: 10.0.1.7:3000, admin, dual1234"):
    #                 grafana = Grafana("", **smallNode)
    #                 v3 >> InvEdge() >> Docker("", **smallNode) >> InvEdge() >> grafana # Grafana("", **smallNode)

    #         #with Cluster("Environment"):
    #         with Cluster("nginx"):
    #             v3 >> InvEdge() >> Nginx("", **smallNode) >> InvEdge() >> Custom(v3_ingest, "./resources/1x1.png", **{"fixedsize": "true", "width": "3.3", "height": "0.4"})


    #         with Cluster("Dual-Link"):
    #             with Cluster("dlapi"):
    #                 grafana >> InvEdge() >> Docker("", **smallNode) >> InvEdge() >> Swift("", **smallNode)
    #             with Cluster("dlsync"):
    #                 grafana >> InvEdge() >> Docker("", **smallNode) >> InvEdge() >> Swift("", **smallNode)
    #             with Cluster("dlscript"):
    #                 grafana >> InvEdge() >> Docker("", **smallNode) >> InvEdge() >> Swift("", **smallNode)
    #             with Cluster("dltemplate"):
    #                 grafana >> InvEdge() >> Docker("", **smallNode) >> InvEdge() >> Swift("", **smallNode)
    #             with Cluster("dlservices"):
    #                 grafana >> InvEdge() >> Docker("", **smallNode) >> InvEdge() >> Php("", **smallNode)
    #             with Cluster("manager"):
    #                 grafana >> InvEdge() >> Docker("", **smallNode) >> InvEdge() >> Custom("", "./resources/www.png", **smallNode)
    #             with Cluster("postgres: 10.0.1.7:5432, root, M10labSpaSS"):
    #                 grafana >> InvEdge() >> Docker("", **smallNode) >> InvEdge() >> Postgresql("", **smallNode)

    #     with Cluster("v3 AWS database backup"):
    #         v3_aws_db_bck = VMLinux("v3-aws-db-bck\n 10.0.1.8 \n dualuser \n 1234")
    #         silentbob >>  InvEdge() >> v3_aws_db_bck >> InvEdge() >> Postgresql("postgres \n Miolabs#Aws-p4assw")

    #     with Cluster("Local v4 Auth for v3 migrations"):
    #         v4_auth = VMLinux("v4-auth\n 192.168.10.12\n 10.0.1.12 \n dualuser \n 1234")
    #         nginx_auth_v4 = Nginx("migration.auth.dual-link.com:8443", **{"fixedsize": "true", "height": "2"})
    #         silentbob >> InvEdge() >> v4_auth >> InvEdge() >> nginx_auth_v4
    #         with Cluster("auth"):
    #             nginx_auth_v4 >> InvEdge() >> Docker("", **smallNode) >> InvEdge() >> Swift("", **smallNode)

    #     with Cluster("Ant Media Server"):
    #         ant = Custom("192.168.10.13:5080\n 10.0.1.9:5080 \n soporte@dual-link \n 14Du@1-link32#", "./resources/Ant-Media.png", **{"fixedsize": "true", "width": "2"} )
    #         silentbob >>  InvEdge() >> VMLinux("ant-01 \n 192.168.10.13\n 10.0.1.9 \n dualuser \n 1234") >> InvEdge() >> ant
    #         with Cluster("Obser"):
    #             with Cluster("Node exporter"):
    #                 Docker("", **smallNode) >> InvEdge() >> Custom("", "./resources/prom-exp.png", **smallNode)
    #             with Cluster("prometheus: 10.0.1.9:9090"):
    #                 Docker("", **smallNode) >> InvEdge() >> Prometheus("", **smallNode)
    #             with Cluster("grafana: 10.0.1.9:3000, admin, dual1234"):
    #                 ant >> InvEdge() >> Docker("", **smallNode) >> InvEdge() >> Grafana("", **smallNode)

    #     with Cluster("Manual Builds"):
    #         build_ws = VMLinux("build-ws\n IP to setup \n dualuser \n 1234")
    #         silentbob >>  InvEdge() >> build_ws >> InvEdge() >> Git("")

    #     with Cluster("Docker Compose Test Cluster"):
    #         docker_test_vm = VMLinux("No existe aún\n 10.0.1.?? \n dualuser \n 1234")
    #         docker_test = Docker("Docker Compose")
    #         silentbob >> InvEdge() >> docker_test_vm >> InvEdge() >> docker_test

    #     with Cluster("Kubernetes Dev Cluster"):
    #         k8s_test_vm = VMLinux("k8s-test\n 10.0.1.11 \n dualuser \n 1234")
    #         k8s_test = K3S("k3s")
    #         silentbob >> InvEdge() >> k8s_test_vm >> InvEdge() >> k8s_test

    #     with Cluster("v4 database test"):
    #         v4_db_test = VMLinux("No existe aún\n 10.0.1.?? \n dualuser \n 1234")
    #         # docker_test >>  InvEdge() >> v4_db_test
    #         # k8s_test >>  InvEdge() >> v4_db_test 
    #         docker_test >>  v4_db_test
    #         k8s_test >>  v4_db_test 
    #         v4_db_test >> InvEdge() >> Postgresql("postgres \n Miolabs#Aws-p4assw")
            
    # #telefonica - unify - silentbob
    # unify - silentbob_pc
    # Vodafone - oficina - silentbob_pc
