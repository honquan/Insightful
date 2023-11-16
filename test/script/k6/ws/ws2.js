import { randomString, randomIntBetween } from 'https://jslib.k6.io/k6-utils/1.2.0/index.js';
import ws from 'k6/ws';
import { check } from 'k6';

export let options = {
    // stages: [
    //     { duration: '10s', target: 4000 },
    //     { duration: '30s', target: 4000 },
    //     { duration: '20s', target: 8000 },
    //     { duration: '30s', target: 8000 },
    //     { duration: '20s', target: 4000 },
    //     { duration: '30s', target: 4000 },
    //     { duration: '10s', target: 0 },
    // ],
    stages: [
        { duration: '30s', target: 8000 },
        // { duration: '7s', target: 6000 },
        // { duration: '7s', target: 6000 },
        // { duration: '7s', target: 4000 },
        // { duration: '5s', target: 0 },
    ],
};

export default function () {
    // var items = Array("index", "contact", "news", "category", "about", "cycle", "sitemap");
    // var item = items[Math.floor(Math.random()*items.length)];
    //
    // // get random userid
    // const userID = randomIntBetween(1, 30000000);
    //
    // // get coordinate array
    // var arrCoordinate = [];
    // for (let i = 0; i < 100; i++) {
    //     const x = randomIntBetween(1, 1980);
    //     const y = randomIntBetween(1, 10000);
    //
    //     arrCoordinate.push([x, y])
    // }
    //
    // var params = ["h", userID, `https://cycle.com/${item}`, arrCoordinate];

    // check fastest running worker ws
    var params = ["h", 1, `https://cycle.com`, []];

    // public websocket server for quick test
    const url = 'ws://localhost:8899/apis/ws/worker-ants';    // apis worker ants
    //const url = 'ws://localhost:8899';    // gnet ws
    //const url = 'ws://localhost:8899/apis/fasthttp-ws/worker-ants';    // fasthttp ws with ants

    const res = ws.connect(url, null, function (socket) {
        /*socket.on('open', function open() {
            //console.log('connected');
            socket.setInterval(function interval() {
                socket.send(JSON.stringify(params));
            }, 1);
        });*/

        while (true) {
            socket.send(JSON.stringify(params));
        }

        // socket.on('message', function message(data) {
        //     console.log('Message received: ', data);
        //     check(data, { 'data is correct': (r) => r && r === text });
        // });

        // socket.on('error', function (e) {
        //     if (e.error() != 'websocket: close sent') {
        //         console.log('An unexpected error occured: ', new Date().toISOString(), e.error());
        //     }
        // });

        //socket.on('close', () => console.log('disconnected'));

        // socket.setTimeout(function () {
        //     console.log('5 seconds passed, closing the socket');
        //     socket.close();
        // }, 5000);
    });

    check(res, { 'status is 101': (r) => r && r.status === 101 });
}