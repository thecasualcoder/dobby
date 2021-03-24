FROM ubuntu:bionic
COPY ./dobby /usr/local/bin/
EXPOSE 4444
RUN useradd -ms /bin/bash dobby
USER dobby
CMD ["dobby", "server", "--bind-address", "0.0.0.0"]
