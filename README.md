newrelic_sphinx
===============

Sphinx monitoring plugin for NewRelic

Installation

If you have not Go compiler in your system:
sudo apt-get install golang

Install dependencies
sudo go get github.com/yunge/sphinx
sudo go get github.com/yvasiyarov/newrelic_platform_go

Get and build agent
git clone https://github.com/yvasiyarov/newrelic_sphinx.git
cd newrelic_sphinx

go build -o sphinx_agent

Run agent in debug mode:
./sphinx_agent --verbose=true --sphinx-host=127.0.0.1 --sphinx-port=9312 --newrelic-license=[your newrelic license key]

In production mode you can run it with nohup
nohup ./sphinx_agent --sphinx-host=127.0.0.1 --sphinx-port=9312 --newrelic-license=[your newrelic license key]

