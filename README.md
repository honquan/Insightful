# Insightful
Insightful

docker build -t "insightful:1.0.1" .

docker run -it --network=host -p 8899:8899 insightful:1.0.1
