insert into goodbuy.product_categories values
('Открытка', 'описание'),
('Стикерпак', 'описание');

insert into goodbuy.products values
('Открытка', 15, 1, 5, 11),
('Стикерпак', 20, 2, 5, 2),
('Еще открытка', 10, 1, 5, 4),
('Новый стикерпак', 11, 2, 5, 3),
('Плоттерная штука', 7, 0, 5, 7),
('Серьги', 50, 0, 5, 100);

insert into goodbuy.positions values
(1, 15, 2, 1),
(2, 20, 1, 1),
(6, 7, 2, 1),
(3, 10, 3, 1),
(3, 10, 2, 1),
(4, 15, 1, 1),
(5, 10, 1, 1),
(6, 30, 1, 1);

insert into goodbuy.receipts values
(date '2023-12-27', 1),
(date '2023-12-27', 1),
(date '2023-12-30', 1),
(date '2023-12-31', 1),
(date '2023-12-31', 1);

insert into goodbuy.positions_in_receipts values
(1, 1),
(2, 1),
(3, 2),
(4, 3),
(5, 4),
(6, 5),
(7, 5),
(8, 5);

insert into goodbuy.markets values
('Супермаркет', '[2023-12-27,2023-12-27]'),
('Гипермаркет', '[2023-12-30, 2023-12-31]');

insert into goodbuy.receipts_on_markets values
(1, 1),
(2, 1),
(3, 2),
(4, 2),
(5, 2);

insert into goodbuy.users values (2, 'Seller', 'd41e98d1eafa6d6011d3a70f1a5b92f0');
insert into goodbuy.users values (3, 'Analytic', 'd41e98d1eafa6d6011d3a70f1a5b92f0');