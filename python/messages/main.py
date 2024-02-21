import asyncio
from src.neo4j_adapter import Neo4jAdapter
from src.consumer import NATSConsumer
from configs.settings import NEO4J_URI, NEO4J_USERNAME, NEO4J_PASSWORD, NATS_URL, NKEYSEED, USE_QUEUE_GROUP


async def main():
    neo4j_adapter = Neo4jAdapter(uri=NEO4J_URI, username=NEO4J_USERNAME, password=NEO4J_PASSWORD)

    # Connect to Neo4j
    await neo4j_adapter.connect()

    consumer = NATSConsumer(
        nats_url=NATS_URL,
        nkeyseed=NKEYSEED,
        subjects=["irc"],
        durable_name="threadr-irc",
        stream_name="messages",
        use_queue_group=USE_QUEUE_GROUP,
        neo4j_adapter=neo4j_adapter
    )
    await consumer.run()

    # Cleanup
    await neo4j_adapter.close()

if __name__ == '__main__':
    asyncio.run(main())
