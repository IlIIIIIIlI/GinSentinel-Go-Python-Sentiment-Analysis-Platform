import grpc
import asyncio
import concurrent.futures
import logging
import os
import sys

# Add generated code directory to path
sys.path.append(os.path.join(os.path.dirname(__file__), 'gen'))

# Import generated protobuf code
from gen.sentiment.v1 import sentiment_pb2
from gen.sentiment.v1 import sentiment_pb2_grpc


# Import the sentiment model
from sentiment_model import analyze_text

# Configure logging
logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s - %(name)s - %(levelname)s - %(message)s'
)
logger = logging.getLogger(__name__)

class SentimentAnalyzerServicer(sentiment_pb2_grpc.SentimentAnalyzerServicer):
    """Provides methods that implement functionality of sentiment analyzer server."""

    async def _analyze(self, text, language):
        """
        Perform sentiment analysis on the given text.
        This is an async method that will call our sentiment model.
        """
        logger.info(f"Analyzing text (length: {len(text)}) in language: {language}")
        try:
            # Call the sentiment analysis model
            result = await analyze_text(text, language)
            logger.info(f"Analysis complete: {result['sentiment']} with score {result['score']}")
            return result
        except Exception as e:
            logger.error(f"Error during text analysis: {e}")
            raise



    def AnalyzeSentiment(self, request, context):
        """
        Implement the AnalyzeSentiment RPC method.
        This method is called by gRPC.
        """
        logger.info(f"Received sentiment analysis request for text: {request.text[:50]}...")

        # Ensure each worker thread has its own event loop
        loop = asyncio.new_event_loop()  # Create a new event loop for the thread
        asyncio.set_event_loop(loop)     # Set it as the current event loop

        # Run the async analysis
        result = loop.run_until_complete(self._analyze(request.text, request.language))

        # Construct the response
        response = sentiment_pb2.SentimentResponse(
            request_id=request.request_id,
            sentiment=result["sentiment"],
            score=result["score"]
        )

        # Add confidence scores
        for key, value in result["confidence_scores"].items():
            response.confidence_scores[key] = value

        # Add keywords
        response.keywords.extend(result["keywords"])

        logger.info(f"Completed analysis for request ID: {request.request_id}")
        return response


    def BatchAnalyzeSentiment(self, request_iterator, context):
        """
        Implement the BatchAnalyzeSentiment RPC method.
        This method handles streaming requests from the client.
        """
        logger.info("Starting batch analysis stream")

        # Ensure each worker thread has its own event loop
        loop = asyncio.new_event_loop()  # Create a new event loop for the thread
        asyncio.set_event_loop(loop)     # Set it as the current event loop

        for request in request_iterator:
            logger.info(f"Processing stream request: {request.request_id}")

            # Analyze the text
            result = loop.run_until_complete(self._analyze(request.text, request.language))

            # Construct and yield the response
            response = sentiment_pb2.SentimentResponse(
                request_id=request.request_id,
                sentiment=result["sentiment"],
                score=result["score"]
            )

            # Add confidence scores
            for key, value in result["confidence_scores"].items():
                response.confidence_scores[key] = value

            # Add keywords
            response.keywords.extend(result["keywords"])

            logger.info(f"Yielding stream response for request ID: {request.request_id}")
            yield response

        logger.info("Batch analysis stream completed")


def serve():
    """Start the gRPC server."""
    port = os.environ.get('GRPC_PORT', '50051')

    # 增加最大工作线程数
    max_workers = 20
    server = grpc.server(concurrent.futures.ThreadPoolExecutor(max_workers=max_workers))
    sentiment_pb2_grpc.add_SentimentAnalyzerServicer_to_server(
        SentimentAnalyzerServicer(), server
    )

    # 确保绑定到所有接口，而不仅仅是localhost
    server_address = f'[::]:{port}'
    server.add_insecure_port(server_address)
    server.start()

    logger.info(f"Server started, listening on {server_address} with {max_workers} workers")

    try:
        server.wait_for_termination()
    except KeyboardInterrupt:
        logger.info("Server stopping due to keyboard interrupt")
        server.stop(0)
        logger.info("Server stopped")

if __name__ == '__main__':
    # 打印系统路径，帮助诊断导入问题
    logger.info(f"Python path: {sys.path}")
    logger.info(f"Starting sentiment analysis gRPC server...")
    serve()