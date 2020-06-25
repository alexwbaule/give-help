INSERT INTO "tags" ("createdat", "name")
VALUES
  ('2020-04-08 00:00:00', 'Construção Cívil'),
  ('2020-04-08 00:00:00', 'Pets'),
  ('2020-04-08 00:00:00', 'Mecânica'),
  ('2020-04-08 00:00:00', 'Alimentação'),
  ('2020-04-10 00:00:00', 'Service Test Category');
INSERT INTO "users" (
    "userid",
    "createdat",
    "name",
    "description",
    "deviceid",
    "allowsharedata",
    "tags",
    "images",
    "giver",
    "taker",
    "url",
    "email",
    "facebook",
    "instagram",
    "twitter",
    "additionaldata",
    "address",
    "city",
    "state",
    "zipcode",
    "country",
    "lat",
    "Lon",
    "lastupdate",
    "registerfrom"
  )
VALUES
  (
    '01E5DEKKFZRKEYCRN6PDXJ8GYZ',
    '2020-04-08 00:00:00',
    'Usuario Da Silva',
    'Nosso querido usuário de testes unitários, agora atualizado',
    '01E5JRTKTTSA6QFK2H3RXFN8SD',
    '1',
    '{"Usuário de testes",TI,"Serviços Gerais"}',
    NULL,
    2.5,
    2.5,
    'usuario.com.br',
    'usuario@email.com',
    'usuario@facebook.com',
    'usuario@instagram.com',
    '@usuario',
    '',
    'Rua da casa do usuário, 777',
    'São Paulo',
    'São Paulo',
    99999000,
    'Brasil',
    -23.5475,
    -46.63611,
    '2020-04-10 19:47:15.195059',
    NULL
  ),
  (
    '01E5DEKKFZRKEYCRN6PDXJ8UUU',
    '2020-04-10 00:00:00',
    'Jose Alterado Pelo Serviço',
    'Nosso querido usuário de testes unitários',
    '01E5JTT75FHE6CP5N187PYZPTP',
    '1',
    '{"Usuário de testes",TI,"Serviços Gerais"}',
    NULL,
    2.5,
    2.5,
    'usuario.com.br',
    'usuario@email.com',
    'usuario@facebook.com',
    'usuario@instagram.com',
    '@usuario',
    '',
    'Rua da casa do usuário, 777',
    'São Paulo',
    'São Paulo',
    99999000,
    'Brasil',
    -23.5475,
    -46.63611,
    '2020-04-10 19:59:36.114563',
    NULL
  );
INSERT INTO "phones" (
    "phoneid",
    "createdat",
    "userid",
    "countrycode",
    "isdefault",
    "phonenumber",
    "region",
    "whatsapp"
  )
VALUES
  (
    119,
    '2020-04-10 00:00:00',
    '01E5DEKKFZRKEYCRN6PDXJ8GYZ',
    '+55',
    '1',
    '88888-8888',
    '11',
    '1'
  ),
  (
    120,
    '2020-04-10 00:00:00',
    '01E5DEKKFZRKEYCRN6PDXJ8GYZ',
    '+55',
    '0',
    '1111-1111',
    '11',
    '0'
  ),
  (
    129,
    '2020-04-10 00:00:00',
    '01E5DEKKFZRKEYCRN6PDXJ8GYZ',
    '+55',
    '1',
    '9999-9999',
    '11',
    '1'
  ),
  (
    130,
    '2020-04-10 00:00:00',
    '01E5DEKKFZRKEYCRN6PDXJ8GYZ',
    '+55',
    '0',
    '1111-1111',
    '11',
    '0'
  );
INSERT INTO "proposals" (
    "proposalid",
    "createdat",
    "userid",
    "side",
    "proposaltype",
    "tags",
    "description",
    "proposalvalidate",
    "lat",
    "lon",
    "range",
    "isactive",
    "lastupdate",
    "areatags",
    "images",
    "exposeuserdata",
    "datatoshare",
    "estimatedvalue",
    "title"
  )
VALUES
  (
    '01E5DEKKFZRKEYCRN6PDXJ8PPP',
    '2020-04-11 00:00:00',
    '01E5DEKKFZRKEYCRN6PDXJ8GYZ',
    'request',
    'product',
    '{Alimentação}',
    'Estou morrendo de fome, adoraria qualquer coisa para comer',
    '2021-06-09 00:00:00',
    -23.5475,
    -46.6361,
    5,
    '1',
    '2020-04-11 01:46:47.906949',
    '{ZL,PENHA,"ZONA LESTE"}',
    '{http://my-domain.com/image1.jpg,http://my-domain.com/image2.jpg,http://my-domain.com/image3.jpg}',
    '1',
    '{Phone,Email,Facebook,Instagram,URL}',
    50.0000,
    'Quero comer'
  ),
  (
    '01E5DEKKFZRKEYCRN6PDXJ8GXX',
    '2020-04-10 00:00:00',
    '01E5DEKKFZRKEYCRN6PDXJ8GYZ',
    'request',
    'product',
    '{Alimentação}',
    'Estou morrendo de fome, adoraria qualquer coisa para comer',
    '2021-06-09 00:00:00',
    -23.5475,
    -46.6361,
    5,
    '1',
    '2020-04-11 02:19:06.427195',
    '{ZL,PENHA,"ZONA LESTE"}',
    '{http://my-domain.com/image1.jpg,http://my-domain.com/image2.jpg,http://my-domain.com/image3.jpg}',
    '1',
    '{Phone,Email,Facebook,Instagram,URL}',
    50.0000,
    'Quero comer'
  );
INSERT INTO "transactions" (
    "transactionid",
    "proposalid",
    "giverid",
    "takerid",
    "createdat",
    "giverrating",
    "giverreviewcomment",
    "takerrating",
    "takerreviewcomment",
    "status",
    "lastupdate"
  )
VALUES
  (
    '01E5DEKKFZRKEYCRN6PDXJ8GZZ',
    '01E5DEKKFZRKEYCRN6PDXJ8GXX',
    '01E5DEKKFZRKEYCRN6PDXJ8GYZ',
    '01E5DEKKFZRKEYCRN6PDXJ8GYZ',
    '2020-04-11 00:00:00',
    4,
    'Cara gente fina, deu tudo certo!',
    5,
    'Foi muito legal, recomendo!',
    'done',
    '2020-04-11 02:03:52.187366'
  ),
  (
    '01E5DEKKFZRKEYCRN6PDXJ8TTT',
    '01E5DEKKFZRKEYCRN6PDXJ8PPP',
    '01E5DEKKFZRKEYCRN6PDXJ8GYZ',
    '01E5DEKKFZRKEYCRN6PDXJ8GYZ',
    '2020-04-11 00:00:00',
    0,
    '',
    0,
    '',
    'done',
    '2020-04-11 02:56:49.296733'
  );
  INSERT INTO TERMS(
    TermID,
    Title,
    Description
) 
VALUES
(
'01E5DEKKFZRKEYCRN6PDXJ8T11',
'Termo de teste 1',
'Lorem Ipsum is simply dummy text of the printing and typesetting industry. 
Lorem Ipsum has been the industrys standard dummy text ever since the 1500s, 
when an unknown printer took a galley of type and scrambled it to make a type specimen book. 
It has survived not only five centuries, but also the leap into electronic typesetting, 
remaining essentially unchanged. It was popularised in the 1960s with the release of Letraset 
sheets containing Lorem Ipsum passages, and more recently with desktop publishing software 
like Aldus PageMaker including versions of Lorem Ipsum.'
);
INSERT INTO TERMS(
    TermID,
    Title,
    Description
) 
VALUES
(
'01E5DEKKFZRKEYCRN6PDXJ8T22',
'Termo de teste 2',
'Lorem Ipsum is simply dummy text of the printing and typesetting industry. 
Lorem Ipsum has been the industrys standard dummy text ever since the 1500s, 
when an unknown printer took a galley of type and scrambled it to make a type specimen book. 
It has survived not only five centuries, but also the leap into electronic typesetting, 
remaining essentially unchanged. It was popularised in the 1960s with the release of Letraset 
sheets containing Lorem Ipsum passages, and more recently with desktop publishing software 
like Aldus PageMaker including versions of Lorem Ipsum.'
);
INSERT INTO TERMS(
    TermID,
    Title,
    Description
) 
VALUES
(
'01E5DEKKFZRKEYCRN6PDXJ8T33',
'Termo de teste 3',
'Lorem Ipsum is simply dummy text of the printing and typesetting industry. 
Lorem Ipsum has been the industrys standard dummy text ever since the 1500s, 
when an unknown printer took a galley of type and scrambled it to make a type specimen book. 
It has survived not only five centuries, but also the leap into electronic typesetting, 
remaining essentially unchanged. It was popularised in the 1960s with the release of Letraset 
sheets containing Lorem Ipsum passages, and more recently with desktop publishing software 
like Aldus PageMaker including versions of Lorem Ipsum.'
);

INSERT INTO TERMS_ACCEPTED 
(
  UserID,
  TermID
) 
VALUES
(
  '01E5DEKKFZRKEYCRN6PDXJ8GYZ',
  '01E5DEKKFZRKEYCRN6PDXJ8T11'
);

INSERT INTO TERMS_ACCEPTED 
(
  UserID,
  TermID
) 
VALUES
(
  '01E5DEKKFZRKEYCRN6PDXJ8GYZ',
  '01E5DEKKFZRKEYCRN6PDXJ8T22'
);

INSERT INTO TERMS_ACCEPTED 
(
  UserID,
  TermID
) 
VALUES
(
  '01E5DEKKFZRKEYCRN6PDXJ8GYZ',
  '01E5DEKKFZRKEYCRN6PDXJ8T33'
);