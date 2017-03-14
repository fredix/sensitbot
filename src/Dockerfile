FROM codenvy/go

# Create the directory where the application will reside
# RUN mkdir /app

# Copy the application files (needed for production)
ADD sensitbot_linux /home/user/application/sensitbot_linux
ADD sensitbot.toml-example /home/user/application/sensitbot.toml
ADD golang.jpg /home/user/application/golang.jpg

# Set the working directory to the app directory
# WORKDIR /app
WORKDIR /home/user/application/

# Expose the application on port 3000
EXPOSE 3000

# Set the entry point of the container to the application executable
# ENTRYPOINT /app/sensitbot_linux
ENTRYPOINT /home/user/application/sensitbot_linux
