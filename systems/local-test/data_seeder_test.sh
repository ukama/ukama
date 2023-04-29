# This script is used to seed the database with test data for the local test environment.
# It is used by the local-test docker-compose file to seed the database with test data.

# REGISTY SYSTEM
# USER SUB-SYSTEM QUERIES
USERS_DB="users"
USER_QUERY="INSERT INTO \"public\".\"users\" (\"id\", \"created_at\", \"updated_at\", \"deleted_at\", \"uuid\", \"name\", \"email\", \"phone\", \"deactivated\") VALUES
(1, '2023-04-29 12:26:22.669305+00', '2023-04-29 12:26:22.669305+00', NULL, '851e0abc-7e91-4206-8c85-498e16f91e67', 'Ukama Root', 'hello@ukama.com', '0000000000', 'f'),
(2, '2023-04-29 12:27:44.53692+00', '2023-04-29 12:27:44.53692+00', NULL, '2b83cb9b-782f-4a75-bb66-0d35538d18d2', 'Salman', 'salman@ukama.com', '3161010101', 'f'),
(3, '2023-04-29 12:27:58.125006+00', '2023-04-29 12:27:58.125006+00', NULL, 'ec4c897e-cc78-43c7-aee3-871a956808c4', 'Brackley', 'brackley@ukama.com', '3162020202', 'f'),
(4, '2023-04-29 12:28:21.218027+00', '2023-04-29 12:28:21.218027+00', NULL, '08a594d7-a292-43cf-9652-54785b03f48f', 'Vishal', 'vishal@ukama.com', '3163030303', 'f'),
(5, '2023-04-29 12:28:47.35359+00', '2023-04-29 12:28:47.35359+00', NULL, '022586f0-2d0f-4b30-967d-2156574fece4', 'Kashif', 'kashif@ukama.com', '3164040404', 'f'),
(6, '2023-04-29 12:29:08.418805+00', '2023-04-29 12:29:08.418805+00', NULL, 'c9647e7a-8967-4978-b512-38a35899f32d', 'Elvis', 'elvis@ukama.com', '3165050505', 'f');"

# ORG SUB-SYSTEM QUERIES
ORG_DB="org"
ORG_QUERY="INSERT INTO \"public\".\"orgs\" (\"id\", \"name\", \"owner\", \"certificate\", \"created_at\", \"updated_at\", \"deactivated\", \"deleted_at\") VALUES
('6c36d521-3bda-4d3f-bddd-375d2e9e2f54', 'ukama', '851e0abc-7e91-4206-8c85-498e16f91e67', '', '2023-04-29 12:26:20.984164+00', '2023-04-29 12:26:20.984164+00', 'f', NULL),
('aac1ed88-2546-4f9c-a808-fb9c4d0ef24b', 'saturn', '08a594d7-a292-43cf-9652-54785b03f48f', 'test_cert', '2023-04-29 12:34:15.365953+00', '2023-04-29 12:34:15.365953+00', 'f', NULL),
('bf184df7-0ce6-4100-a9c6-497c181b87cf', 'milky-way', '022586f0-2d0f-4b30-967d-2156574fece4', 'test_cert', '2023-04-29 12:31:27.092514+00', '2023-04-29 12:31:27.092514+00', 'f', NULL);"

USERS_IN_ORG_QUERY="INSERT INTO \"public\".\"users\" (\"id\", \"uuid\", \"deactivated\", \"deleted_at\") VALUES
(1, '851e0abc-7e91-4206-8c85-498e16f91e67', 'f', NULL),
(2, '2b83cb9b-782f-4a75-bb66-0d35538d18d2', 'f', NULL),
(3, 'ec4c897e-cc78-43c7-aee3-871a956808c4', 'f', NULL),
(4, '08a594d7-a292-43cf-9652-54785b03f48f', 'f', NULL),
(5, '022586f0-2d0f-4b30-967d-2156574fece4', 'f', NULL),
(6, 'c9647e7a-8967-4978-b512-38a35899f32d', 'f', NULL);"

ORG_USERS_QUERY="INSERT INTO \"public\".\"org_users\" (\"org_id\", \"user_id\", \"uuid\", \"deactivated\", \"created_at\", \"deleted_at\") VALUES
('6c36d521-3bda-4d3f-bddd-375d2e9e2f54', 1, '851e0abc-7e91-4206-8c85-498e16f91e67', 'f', '2023-04-29 12:26:21.05832+00', NULL),
('6c36d521-3bda-4d3f-bddd-375d2e9e2f54', 2, '2b83cb9b-782f-4a75-bb66-0d35538d18d2', 'f', '2023-04-29 12:27:44.606257+00', NULL),
('6c36d521-3bda-4d3f-bddd-375d2e9e2f54', 3, 'ec4c897e-cc78-43c7-aee3-871a956808c4', 'f', '2023-04-29 12:27:58.130211+00', NULL),
('6c36d521-3bda-4d3f-bddd-375d2e9e2f54', 4, '08a594d7-a292-43cf-9652-54785b03f48f', 'f', '2023-04-29 12:28:21.227499+00', NULL),
('6c36d521-3bda-4d3f-bddd-375d2e9e2f54', 5, '022586f0-2d0f-4b30-967d-2156574fece4', 'f', '2023-04-29 12:28:47.362703+00', NULL),
('6c36d521-3bda-4d3f-bddd-375d2e9e2f54', 6, 'c9647e7a-8967-4978-b512-38a35899f32d', 'f', '2023-04-29 12:29:08.431079+00', NULL),
('aac1ed88-2546-4f9c-a808-fb9c4d0ef24b', 4, '08a594d7-a292-43cf-9652-54785b03f48f', 'f', '2023-04-29 12:34:15.364808+00', NULL),
('aac1ed88-2546-4f9c-a808-fb9c4d0ef24b', 5, '022586f0-2d0f-4b30-967d-2156574fece4', 'f', '2023-04-29 12:34:38.620427+00', NULL),
('aac1ed88-2546-4f9c-a808-fb9c4d0ef24b', 6, 'c9647e7a-8967-4978-b512-38a35899f32d', 'f', '2023-04-29 12:34:50.209656+00', NULL),
('bf184df7-0ce6-4100-a9c6-497c181b87cf', 2, '2b83cb9b-782f-4a75-bb66-0d35538d18d2', 'f', '2023-04-29 12:32:54.858956+00', NULL),
('bf184df7-0ce6-4100-a9c6-497c181b87cf', 3, 'ec4c897e-cc78-43c7-aee3-871a956808c4', 'f', '2023-04-29 12:33:13.822997+00', NULL),
('bf184df7-0ce6-4100-a9c6-497c181b87cf', 5, '022586f0-2d0f-4b30-967d-2156574fece4', 'f', '2023-04-29 12:31:27.091411+00', NULL);"

# NETWORK SUB SYSTEM QUERIES
NETWORK_DB="network"
NETWORKS_QUERY="INSERT INTO \"public\".\"networks\" (\"id\", \"name\", \"org_id\", \"deactivated\", \"created_at\", \"updated_at\", \"deleted_at\") VALUES
('06455edb-d33b-49ba-b8ed-589cf718047a', 'mesh-network', 'bf184df7-0ce6-4100-a9c6-497c181b87cf', 'f', '2023-04-29 12:38:06.972679+00', '2023-04-29 12:38:06.972679+00', NULL),
('b884485f-cb43-44b1-be57-0b777b154ff2', 'saturn-network', 'aac1ed88-2546-4f9c-a808-fb9c4d0ef24b', 'f', '2023-04-29 12:38:28.660736+00', '2023-04-29 12:38:28.660736+00', NULL);"

NETWORK_ORGS_QUERY="INSERT INTO \"public\".\"orgs\" (\"id\", \"name\", \"deactivated\", \"created_at\", \"updated_at\", \"deleted_at\") VALUES
('aac1ed88-2546-4f9c-a808-fb9c4d0ef24b', 'saturn', 'f', '2023-04-29 12:38:28.658766+00', '2023-04-29 12:38:28.658766+00', NULL),
('bf184df7-0ce6-4100-a9c6-497c181b87cf', 'milky-way', 'f', '2023-04-29 12:38:06.96378+00', '2023-04-29 12:38:06.96378+00', NULL);"

SITES_QUERY="INSERT INTO \"public\".\"sites\" (\"id\", \"name\", \"network_id\", \"deactivated\", \"created_at\", \"updated_at\", \"deleted_at\") VALUES
('01cf8020-5d92-4364-a8e5-546208d8859a', 's3-site', '06455edb-d33b-49ba-b8ed-589cf718047a', 'f', '2023-04-29 12:41:05.085422+00', '2023-04-29 12:41:05.085422+00', NULL),
('2dfdd147-0738-43e7-946b-4dd37f99ee88', 's1-site', 'b884485f-cb43-44b1-be57-0b777b154ff2', 'f', '2023-04-29 12:40:14.29955+00', '2023-04-29 12:40:14.29955+00', NULL),
('63dbde62-2769-4505-8fff-c205b71e3fbc', 's2-site', '06455edb-d33b-49ba-b8ed-589cf718047a', 'f', '2023-04-29 12:40:59.175017+00', '2023-04-29 12:40:59.175017+00', NULL),
('c1219d8b-3e8f-4b4b-9784-f9bf1870ec02', 's2-site', 'b884485f-cb43-44b1-be57-0b777b154ff2', 'f', '2023-04-29 12:40:21.23524+00', '2023-04-29 12:40:21.23524+00', NULL),
('df656312-89d1-49ca-9091-45891db71010', 's3-site', 'b884485f-cb43-44b1-be57-0b777b154ff2', 'f', '2023-04-29 12:40:26.330752+00', '2023-04-29 12:40:26.330752+00', NULL),
('e8a93d4d-704b-4507-af58-a16b9a4657b1', 's1-site', '06455edb-d33b-49ba-b8ed-589cf718047a', 'f', '2023-04-29 12:40:55.797713+00', '2023-04-29 12:40:55.797713+00', NULL);"

#BASE RATE SUB SYSTEM QUERIES
BASERATE_DB="baserate"
BASERATE_QUERY="INSERT INTO \"public\".\"base_rates\" (\"id\", \"created_at\", \"updated_at\", \"deleted_at\", \"uuid\", \"country\", \"provider\", \"vpmn\", \"imsi\", \"sms_mo\", \"sms_mt\", \"data\", \"x2g\", \"x3g\", \"x5g\", \"lte\", \"lte_m\", \"apn\", \"effective_at\", \"end_at\", \"sim_type\", \"currency\") VALUES
(1, '2023-04-29 13:15:57.357091+00', '2023-04-29 13:15:57.357091+00', NULL, 'db7a5d5c-7b18-46ff-839d-c9bf15fa90f3', 'The lunar maria', 'ABC Tel', 'TLM', 1, 0.1, 0.1, 0.1, 't', 'f', 'f', 'f', 'f', 'ukama', '2023-05-09 11:00:00+00', '2024-11-06 23:00:00+00', 1, 'Dollar'),
(2, '2023-04-29 13:15:57.379126+00', '2023-04-29 13:15:57.379126+00', NULL, '100a0372-8f08-45ad-8f51-2c558d180709', 'The lunar maria', 'Light Tel', 'TLM', 1, 0.2, 0.1, 0.2, 't', 'f', 'f', 't', 'f', 'ukama', '2023-05-09 11:00:00+00', '2024-11-06 23:00:00+00', 1, 'Dollar'),
(3, '2023-04-29 13:15:57.382355+00', '2023-04-29 13:15:57.382355+00', NULL, 'd377b57a-144c-4472-b43f-eb6bc585cb45', 'The lunar maria', 'Eagle Tel', 'TLM', 1, 0.1, 0.1, 0.3, 't', 't', 'f', 't', 'f', 'ukama', '2023-05-09 11:00:00+00', '2024-11-06 23:00:00+00', 1, 'Dollar'),
(4, '2023-04-29 13:15:57.384781+00', '2023-04-29 13:15:57.384781+00', NULL, 'cccbb722-a273-4f8e-819a-6bce3be85ed2', 'Montes Apenninus', 'Power Tel', 'TMA', 1, 0.2, 0.1, 0.2, 't', 't', 'f', 't', 'f', 'saturn', '2023-05-09 11:00:00+00', '2024-11-06 23:00:00+00', 1, 'Dollar'),
(5, '2023-04-29 13:15:57.387976+00', '2023-04-29 13:15:57.387976+00', NULL, 'e96aa51b-c37c-4e80-99be-51c2eb97819d', 'Montes Apenninus', '2D Tel', 'TMA', 1, 0.1, 0.1, 0.1, 't', 't', 'f', 't', 'f', 'saturn', '2023-05-09 11:00:00+00', '2024-11-06 23:00:00+00', 1, 'Dollar'),
(6, '2023-04-29 13:15:57.390552+00', '2023-04-29 13:15:57.390552+00', NULL, '6ff5e957-e2f7-48e6-8948-912db68a0675', 'Tycho crater', 'Multi Tel', 'TTC', 1, 0.1, 0.1, 0.4, 't', 't', 'f', 't', 'f', 'saturn', '2023-05-09 11:00:00+00', '2024-11-06 23:00:00+00', 1, 'Dollar'),
(7, '2023-04-29 13:15:57.393123+00', '2023-04-29 13:15:57.393123+00', NULL, '8042fea5-dc00-4da2-8cc5-72a4d545962f', 'Tycho crater', 'Connect Tel', 'TTC', 1, 0.1, 0.1, 0.1, 't', 't', 'f', 't', 'f', 'milky-way', '2023-05-09 11:00:00+00', '2024-11-06 23:00:00+00', 1, 'Dollar'),
(8, '2023-04-29 13:15:57.395993+00', '2023-04-29 13:15:57.395993+00', NULL, '8dcf0d22-31ee-42a0-a64c-54f8bf721149', 'Tycho crater', 'OWS Tel', 'TTC', 1, 0.2, 0.1, 0.5, 't', 't', 'f', 't', 'f', 'milky-way', '2023-05-09 11:00:00+00', '2024-11-06 23:00:00+00', 1, 'Dollar');"

RATE_DB="rate"
MARKUP_DEFAULT_QUERY="INSERT INTO \"public\".\"default_markups\" (\"id\", \"created_at\", \"updated_at\", \"deleted_at\", \"markup\") VALUES
(1, '2023-04-29 14:42:52.880065+00', '2023-04-29 14:42:52.880065+00', NULL, 8);"

MARKUPS_QUERY="INSERT INTO \"public\".\"markups\" (\"id\", \"created_at\", \"updated_at\", \"deleted_at\", \"owner_id\", \"markup\") VALUES
(1, '2023-04-29 14:44:23.738064+00', '2023-04-29 14:44:23.738064+00', NULL, '2b83cb9b-782f-4a75-bb66-0d35538d18d2', 10),
(2, '2023-04-29 14:45:05.274174+00', '2023-04-29 14:45:05.274174+00', NULL, 'ec4c897e-cc78-43c7-aee3-871a956808c4', 8),
(3, '2023-04-29 14:45:27.771591+00', '2023-04-29 14:45:27.771591+00', NULL, 'c9647e7a-8967-4978-b512-38a35899f32d', 5),
(4, '2023-04-29 14:45:38.330433+00', '2023-04-29 14:45:38.330433+00', NULL, '022586f0-2d0f-4b30-967d-2156574fece4', 11);"

PACKAGE_DB="package"
PACKAGES_QUERY="INSERT INTO \"public\".\"packages\" (\"id\", \"created_at\", \"updated_at\", \"deleted_at\", \"uuid\", \"owner_id\", \"name\", \"sim_type\", \"org_id\", \"active\", \"duration\", \"sms_volume\", \"data_volume\", \"voice_volume\", \"type\", \"data_units\", \"voice_units\", \"message_units\", \"flatrate\", \"currency\", \"from\", \"to\", \"country\", \"provider\") VALUES
(1, '2023-04-29 14:58:39.503086+00', '2023-04-29 14:58:39.503086+00', NULL, '0e175fc4-2710-4cf7-92f9-c815379252c5', '022586f0-2d0f-4b30-967d-2156574fece4', 'Monthly-Data', 1, 'bf184df7-0ce6-4100-a9c6-497c181b87cf', 't', 0, 0, 2048, 0, 2, 3, 1, 0, 'f', 'Dollar', '2024-01-01 00:00:00+00', '2024-05-01 00:00:00+00', 'Tycho crater', 'Connect Tel'),
(2, '2023-04-29 14:59:58.394306+00', '2023-04-29 14:59:58.394306+00', NULL, '15e2d1b1-6194-4fee-9c32-7c1e8b9184b7', '022586f0-2d0f-4b30-967d-2156574fece4', 'Monthly-Data', 1, 'bf184df7-0ce6-4100-a9c6-497c181b87cf', 't', 0, 0, 1024, 0, 2, 3, 1, 0, 'f', 'Dollar', '2024-01-01 00:00:00+00', '2024-03-01 00:00:00+00', 'Tycho crater', 'Connect Tel'),
(3, '2023-04-29 15:03:57.966069+00', '2023-04-29 15:03:57.966069+00', NULL, '881daafe-2750-45c8-bf71-14c94fb5542d', '08a594d7-a292-43cf-9652-54785b03f48f', 'Monthly-Data', 1, 'aac1ed88-2546-4f9c-a808-fb9c4d0ef24b', 't', 0, 0, 1024, 0, 2, 3, 1, 0, 'f', 'Dollar', '2024-01-01 00:00:00+00', '2024-03-01 00:00:00+00', 'Tycho crater', 'Multi Tel'),
(4, '2023-04-29 15:04:30.426953+00', '2023-04-29 15:04:30.426953+00', NULL, '9bd5dd1f-b8cd-4277-ac46-3695f51596ff', '08a594d7-a292-43cf-9652-54785b03f48f', 'Monthly-Data', 1, 'aac1ed88-2546-4f9c-a808-fb9c4d0ef24b', 't', 0, 0, 1500, 0, 2, 3, 1, 0, 'f', 'Dollar', '2024-01-01 00:00:00+00', '2024-04-01 00:00:00+00', 'Montes Apenninus', '2D Tel'),
(5, '2023-04-29 15:05:20.791156+00', '2023-04-29 15:05:20.791156+00', NULL, 'a3e18b0e-d33d-4cc4-8a31-afd1e0f25b14', '08a594d7-a292-43cf-9652-54785b03f48f', 'Monthly-Data', 1, 'aac1ed88-2546-4f9c-a808-fb9c4d0ef24b', 't', 0, 0, 2000, 0, 2, 3, 1, 0, 'f', 'Dollar', '2024-01-01 00:00:00+00', '2024-05-01 00:00:00+00', 'Montes Apenninus', 'Power Tel');"

PACKAGE_DETAILS_QUERY="INSERT INTO \"public\".\"package_details\" (\"id\", \"created_at\", \"updated_at\", \"deleted_at\", \"package_id\", \"dlbr\", \"ulbr\", \"apn\") VALUES
(1, '2023-04-29 14:58:39.53668+00', '2023-04-29 14:58:39.53668+00', NULL, '0e175fc4-2710-4cf7-92f9-c815379252c5', 10240000, 10240000, 'milky-way'),
(2, '2023-04-29 14:59:58.403962+00', '2023-04-29 14:59:58.403962+00', NULL, '15e2d1b1-6194-4fee-9c32-7c1e8b9184b7', 10240000, 10240000, 'milky-way'),
(3, '2023-04-29 15:03:57.977475+00', '2023-04-29 15:03:57.977475+00', NULL, '881daafe-2750-45c8-bf71-14c94fb5542d', 10240000, 10240000, 'saturn'),
(4, '2023-04-29 15:04:30.433425+00', '2023-04-29 15:04:30.433425+00', NULL, '9bd5dd1f-b8cd-4277-ac46-3695f51596ff', 10240000, 10240000, 'saturn'),
(5, '2023-04-29 15:05:20.797005+00', '2023-04-29 15:05:20.797005+00', NULL, 'a3e18b0e-d33d-4cc4-8a31-afd1e0f25b14', 10240000, 10240000, 'saturn');"

PACKAGE_MARKUPS_QUERY="INSERT INTO \"public\".\"package_markups\" (\"id\", \"created_at\", \"updated_at\", \"deleted_at\", \"package_id\", \"base_rate_id\", \"markup\") VALUES
(1, '2023-04-29 14:58:39.533716+00', '2023-04-29 14:58:39.533716+00', NULL, '0e175fc4-2710-4cf7-92f9-c815379252c5', '8042fea5-dc00-4da2-8cc5-72a4d545962f', 5),
(2, '2023-04-29 14:59:58.4025+00', '2023-04-29 14:59:58.4025+00', NULL, '15e2d1b1-6194-4fee-9c32-7c1e8b9184b7', '8dcf0d22-31ee-42a0-a64c-54f8bf721149', 5),
(3, '2023-04-29 15:03:57.976087+00', '2023-04-29 15:03:57.976087+00', NULL, '881daafe-2750-45c8-bf71-14c94fb5542d', '6ff5e957-e2f7-48e6-8948-912db68a0675', 5),
(4, '2023-04-29 15:04:30.432188+00', '2023-04-29 15:04:30.432188+00', NULL, '9bd5dd1f-b8cd-4277-ac46-3695f51596ff', 'e96aa51b-c37c-4e80-99be-51c2eb97819d', 5),
(5, '2023-04-29 15:05:20.795435+00', '2023-04-29 15:05:20.795435+00', NULL, 'a3e18b0e-d33d-4cc4-8a31-afd1e0f25b14', 'cccbb722-a273-4f8e-819a-6bce3be85ed2', 5);"

PACKAGE_RATES_QUERY="INSERT INTO \"public\".\"package_rates\" (\"id\", \"created_at\", \"updated_at\", \"deleted_at\", \"package_id\", \"amount\", \"sms_mo\", \"sms_mt\", \"data\") VALUES
(1, '2023-04-29 14:58:39.520023+00', '2023-04-29 14:58:39.520023+00', NULL, '0e175fc4-2710-4cf7-92f9-c815379252c5', 227.328, 0, 0, 0.111),
(2, '2023-04-29 14:59:58.399358+00', '2023-04-29 14:59:58.399358+00', NULL, '15e2d1b1-6194-4fee-9c32-7c1e8b9184b7', 113.664, 0, 0, 0.111),
(3, '2023-04-29 15:03:57.972095+00', '2023-04-29 15:03:57.972095+00', NULL, '881daafe-2750-45c8-bf71-14c94fb5542d', 442.36800000000005, 0, 0, 0.43200000000000005),
(4, '2023-04-29 15:04:30.429571+00', '2023-04-29 15:04:30.429571+00', NULL, '9bd5dd1f-b8cd-4277-ac46-3695f51596ff', 162.00000000000003, 0, 0, 0.10800000000000001),
(5, '2023-04-29 15:05:20.793816+00', '2023-04-29 15:05:20.793816+00', NULL, 'a3e18b0e-d33d-4cc4-8a31-afd1e0f25b14', 432.00000000000006, 0, 0, 0.21600000000000003);"

export ORG_QUERY
export USERS_IN_ORG_QUERY
export ORG_USERS_QUERY
export USER_QUERY
export NETWORKS_QUERY
export NETWORK_ORGS_QUERY
export SITES_QUERY
export USERS_DB
export ORGS_DB
export NETWORKS_DB

export BASERATE_DB
export BASERATE_QUERY
export RATE_DB
export MARKUP_DEFAULT_QUERY
export MARKUPS_QUERY
export PACKAGE_DB
export PACKAGES_QUERY
export PACKAGE_DETAILS_QUERY
export PACKAGE_MARKUPS_QUERY
export PACKAGE_RATES_QUERY