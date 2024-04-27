document.addEventListener('DOMContentLoaded', () => {
    const url = 'https://api.coingecko.com/api/v3/coins/markets?vs_currency=usd&order=market_cap_desc&per_page=250&page=1';
    fetchCryptoData(url);
});

function fetchCryptoData(url) {
    fetch(url)
        .then(response)
        .then(createTable)
        .catch(handleFetchError);
}

function response(response) {
    if (!response.ok) {
        throw new Error(`error: ${response.status}`);
    }
    return response.json();
}

function createTable(data) {
    const table = document.querySelector('#cryptoTable tbody');
    data.forEach((crypto, i) => {
        const row = createTableRow(crypto, i);
        table.appendChild(row);
    });
}

function createTableRow(coin, index) {
    const row = document.createElement('tr');
    row.appendChild(createTableCell(coin.id));
    row.appendChild(createTableCell(coin.symbol));
    row.appendChild(createTableCell(coin.name));

    if (index < 5) {
        row.classList.add('blue-background');
    }

    if (coin.symbol === 'usdt') {
        row.classList.add('green-background');
    }

    return row;
}

function createTableCell(text) {
    const cell = document.createElement('td');
    cell.textContent = text;
    return cell;
}

function handleFetchError(error) {
    console.error('Fail:', error);
    alert('Failed to load');
}

