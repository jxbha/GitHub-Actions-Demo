FROM alpine:3.20
WORKDIR /app
COPY mana .
RUN chmod +x mana
EXPOSE 4040
CMD ["./mana"]
