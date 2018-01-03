

mkdir /opt/config/disco-bit/
git clone https://github.com/tzapu/disco-bit
cd disco-bit
docker build -t disco-bit . 

docker run \
-d --restart always \
--name disco-bit \
-e D_TOKEN="Bot tOkEn.From.DiSc0rd" \
-v /opt/config/disco-bit/:/go/src/app/config/ \
disco-bit:latest -v
