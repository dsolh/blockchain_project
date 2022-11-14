pushd network

./startnetwork.sh
sleep 60

./createchannel.sh
sleep 30

./setAnchorPeerUpdate.sh
sleep 30

popd