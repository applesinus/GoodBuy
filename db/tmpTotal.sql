select sum(p.count * p.cost) as total
from goodbuy.positions p
    join goodbuy.positions_in_receipts pr on p.id = pr.position
    join goodbuy.receipts r on pr.receipt = r.id
where r.status = 1 and p.status = 1;