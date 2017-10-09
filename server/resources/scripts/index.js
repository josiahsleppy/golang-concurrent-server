; (function () {
    const responseDiv = document.querySelector('.js-response-div');

    document.querySelector('.js-send-requests')
        .addEventListener('click', () => {
            responseDiv.innerHTML = '';
            let userInput = document.querySelector('.js-input-box').value;
            let requestTimes = document.querySelector('.js-request-options').value;
            if (requestTimes > 1000) {
                requestTimes = 1000;
            }
            makeRequests('/api/collatz?num=' + userInput, requestTimes).then(resp => {
                return resp.text();
            }).then(text => {
                responseDiv.innerHTML = text;
            }).catch(err => {
                console.log('What went wrong:', err);
            });
        });

    async function makeRequests(uri, times) {
        let value;
        let uris = [];

        for (let i = 0; i < times; i++) {
            uris.push(uri);
        }

        try {
            for (let request of uris.map(getResource)) {
                let result = await request;
                console.log(result.statusText);
                value = value || result;
            }
        } catch (e) {
            console.error('Major error:', e);
        }

        return value;
    }

    function getResource(uri) {
        try {
            return fetch(uri);
        } catch (e) {
            console.log('Could not fetch resource:', e);
        }
    }

})();
