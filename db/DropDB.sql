select * from goodbuy.get_top_N_products_by_sales(500, 2);
select * from goodbuy.get_top_N_products_by_profit(500, 3);
select * from goodbuy.get_n_popular_products_on_markets(2);
select * from goodbuy.positions;

drop trigger if exists update_products_category_trigger_on_delete on goodbuy.product_categories;
drop function if exists goodbuy.update_products_category;
drop trigger if exists check_default_cost_trigger_on_insert on goodbuy.products;
drop function if exists goodbuy.check_default_cost();

drop function if exists goodbuy.get_top_N_products_by_sales;
drop function if exists goodbuy.get_top_N_products_by_profit;
drop function if exists goodbuy.get_N_popular_products_on_markets;

drop procedure if exists goodbuy.edit_product;
drop procedure if exists goodbuy.add_product;

drop view if exists goodbuy.get_products;
drop view if exists goodbuy.get_detailed_receipts;

drop table if exists goodbuy.receipts_on_markets;
drop table if exists goodbuy.positions_in_receipts;
drop table if exists goodbuy.receipts;
drop table if exists goodbuy.positions;
drop table if exists goodbuy.products;
drop table if exists goodbuy.statuses;
drop table if exists goodbuy.markets;
drop table if exists goodbuy.product_categories;
drop table if exists goodbuy.users;
drop table if exists goodbuy.roles;

drop schema if exists goodbuy;