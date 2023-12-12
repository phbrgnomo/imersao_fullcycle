# Build the image
FROM golang:latest

# Upgrade existing packages, Install additional packages, and clean cache 
RUN apt-get update && \
    apt-get full-upgrade -y &&\
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*

# Install testing tools
RUN go install -v golang.org/x/tools/gopls@latest

ENV USER_NAME=${USER_NAME}
ENV USER_ID=${USER_ID}
ENV GROUP_ID=${GROUP_ID}

# Create a non-root user use the same as the host user
ARG USER_NAME
ARG USER_ID
ARG GROUP_ID

# Add the user to the staff group
RUN groupadd -g $GROUP_ID $USER_NAME && \
    useradd -u $USER_ID -g $GROUP_ID -m -s /bin/bash $USER_NAME

# Give group write permissions on the go directory to be able to install packages
RUN chmod -R 777 /go/pkg

# Set the working directory. NOTE: Build.sh script will only work if you directly add the username, instead of the variable
WORKDIR /home/$USER_NAME/code_env

# Set correct permissions
RUN chmod -R 777 /home/$USER_NAME/code_env

# Switch to the non-root user
USER $USER_NAME

# Set the default command to keep the container running
CMD ["sh", "-c", "while true; do echo 'Container is running...'; sleep 10; done"]

