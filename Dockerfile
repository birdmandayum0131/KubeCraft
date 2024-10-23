FROM python:3.12.6-alpine
LABEL author="Bird"
# Create a non-root user and set the home directory
RUN addgroup -S Servers && adduser -S fastapi -G Servers
# Switch to the non-root user
USER fastapi
# Create the server directory
RUN mkdir /home/fastapi/MinecraftBridge
COPY ./ /home/fastapi/MinecraftBridge
RUN pip install -r /home/fastapi/MinecraftBridge/requirements.txt
# Copy server file and eula to the image for building
WORKDIR /home/fastapi/MinecraftBridge
CMD ["python", "./main.py"]