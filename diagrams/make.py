from diagrams.onprem.queue import Kafka
from diagrams.onprem.compute import Server
from diagrams.onprem.database import PostgreSQL
from diagrams import Cluster, Diagram, Edge

with Diagram("Monitoring pipeline", show=False):
    with Cluster("Monitoring Agents"):
        agents = [
            Server("Agent1"),
            Server("Agent2"),
            Server("Agent3")
        ]
    with Cluster("Targets"):
        targets = [
            Server("Target1") << Edge(color="darkgreen", label="probe") << agents[0],
            Server("Target2") << Edge(color="darkred", label="probe") << agents[1],
            Server("Target3") << Edge(color="darkgreen", label="probe") << agents[2]
        ]
    with Cluster("agent-probe-results.v1.json"):
        topic = Kafka("")
        agents >> Edge(label="push", reverse=False, forward=True) >> topic

    with Cluster("Results processor"):
        processor = Server("Results processor")
        processor << Edge(label="consume") << topic

    database = PostgreSQL("database.probe_results")
    processor >> database
