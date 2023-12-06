document.querySelector('button').onclick = myClick;

function myClick() {
    event.preventDefault(); // Предотвращение действия по умолчанию
    var xhr = new XMLHttpRequest();
    var inputValue = document.getElementById("IDInput").value;
    xhr.open("GET", "http://localhost:8080/api/orders/" + inputValue, true);
    xhr.setRequestHeader("Access-Control-Allow-Origin", "http://localhost:8080");

    xhr.onreadystatechange = function () {
        var tableOrders = document.getElementById('orders');
        var tableDelivery = document.getElementById('delivery');
        var tablePayment = document.getElementById('payment');
        var tableItems = document.getElementById('items');
        if (xhr.readyState === 4 && xhr.status === 200) {
            var response = JSON.parse(xhr.responseText);

            // Найдите таблицу в HTML

// Очистить существующие строки таблицы (кроме заголовка)
            while (tableOrders.rows.length > 1) {
                tableOrders.deleteRow(1);
            }
            while (tableDelivery.rows.length > 1) {
                tableDelivery.deleteRow(1);
            }
            while (tablePayment.rows.length > 1) {
                tablePayment.deleteRow(1);
            }
            while (tableItems.rows.length > 1) {
                tableItems.deleteRow(1);
            }

// Создать новую строку
            var row = tableOrders.insertRow();
// Добавить ячейки с данными в строку
            for (var key in response) {
                if (key==='delivery'){
                    var rowDelivery = tableDelivery.insertRow();
                    console.log(response.delivery);
                    for (var keyDelivery in response.delivery){
                        console.log(keyDelivery);
                        if (keyDelivery==='delivery_uuid'){
                            continue;
                        }
                        if (response.delivery.hasOwnProperty(keyDelivery)) {
                            var cell = rowDelivery.insertCell();
                            cell.innerHTML = response.delivery[keyDelivery];
                        }
                    }
                    continue;
                }
                if (key==='payment'){
                    var rowPayment = tablePayment.insertRow();
                    for (var keyPayment in response.payment){
                        if (keyPayment==='payment_uuid'){
                            continue;
                        }
                        if (response.payment.hasOwnProperty(keyPayment)) {
                            var cell = rowPayment.insertCell();
                            cell.innerHTML = response.payment[keyPayment];
                        }
                    }
                    continue;
                }
                if (key==='items'){
                    console.log(response.items)
                    for (var item in response.items){
                        console.log(response.items[item])
                        var rowItems = tableItems.insertRow();
                        for (var keyItem in response.items[item]){
                            if (keyItem==='item_uuid'){
                                continue;
                            }
                            if (response.items[item].hasOwnProperty(keyItem)) {
                                var cell = rowItems.insertCell();
                                cell.innerHTML = response.items[item][keyItem];
                            }
                        }
                    }
                    //обработать таблицу items
                    continue;
                }
                if (response.hasOwnProperty(key)) {
                    var cell = row.insertCell();
                    cell.innerHTML = response[key];
                }
            }
        } else {
            while (tableOrders.rows.length > 1) {
                tableOrders.deleteRow(1);
            }
            while (tableDelivery.rows.length > 1) {
                tableDelivery.deleteRow(1);
            }
            while (tablePayment.rows.length > 1) {
                tablePayment.deleteRow(1);
            }
            while (tableItems.rows.length > 1) {
                tableItems.deleteRow(1);
            }
            console.error('Произошла ошибка при выполнении запроса.');
        }
    };

    xhr.onerror = function () {
        console.error('Произошла ошибка при выполнении запроса.');
    };

    xhr.send();
}