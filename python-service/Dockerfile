# 使用官方 Python 运行时作为父镜像
FROM python:3.9-slim

# 设置工作目录
WORKDIR /app

# 安装必要的依赖
RUN pip install --no-cache-dir grpcio grpcio-tools

# 复制 proto 文件到容器中
COPY proto /app/proto

# 在容器中创建目标目录
RUN mkdir -p /app/python-service/gen

# 生成 Python 代码
RUN python -m grpc_tools.protoc --proto_path=./proto --python_out=./python-service/gen --grpc_python_out=./python-service/gen ./proto/sentiment/v1/sentiment.proto

# 复制应用程序代码
COPY python-service/server.py /app/python-service/
COPY python-service/sentiment_model.py /app/python-service/

# 设置 PYTHONPATH 环境变量，指向正确的路径
ENV PYTHONPATH="/app/python-service/gen:$PYTHONPATH"

# 暴露端口
EXPOSE 50051

# 设置容器启动命令
CMD ["python", "/app/python-service/server.py"]
