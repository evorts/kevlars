-- migrate:up
insert into clients(id, name, secret, expired_at)
values
(1, 'test_client_1', 'D69430B64AA8B15432F7268A3735391634A0E79C4EB07DA29F029865E43A6799B036A144289B9C60DD318A3B371B7A516C4356F097B7D2C2FF13BB4386714CDB', null),
(2, 'test_client_2', '03BE957A00B7ABFB50DFEE7CEF9577FADE6DE469B93CF7CD347AABD00BB5316DDA6289AB7B37471779D010B19C26E5035D24DCFFFEA24C505E4840BD88063D43', current_timestamp + interval '2' hour),
(3, 'test_client_3', 'A2C24F9BF70E456174EB921F9EED6CFD12BBA784F84E054C27D24E57DDFDEFF03230ABD1E6B9F89825CE8BC189C3EE1DC077756E62EC1D4D197AD33C7B87F540', current_timestamp - interval '2' hour)
;

insert into client_scopes(client_id,resource,scopes)
values
(1,'/res/a', null),
(2,'/res/a', '["write","read"]'::jsonb),
(2,'/res/b', '["write"]'::jsonb),
(2,'/res/c', '["read"]'::jsonb),
(3,'/res/a', '["write","read"]'::jsonb)
;
-- migrate:down
truncate clients;
