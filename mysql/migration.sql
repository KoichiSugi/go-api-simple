CREATE DATABASE IGD;

USE IGD;

DROP TABLE IF EXISTS `employee`;

CREATE TABLE `employee` (
    id VARCHAR(50) NOT NULL PRIMARY KEY,
    first_name VARCHAR(50) NOT NULL,
    middle_name VARCHAR(50) NOT NULL,
    last_name VARCHAR(50) NOT NULL,
    gender VARCHAR(10) NOT NULL,
    salary INTEGER NOT NULL,
    dob DATE,
    email VARCHAR(50) UNIQUE NOT NULL,
    phone VARCHAR(10) NOT NULL,
    state VARCHAR(10) NOT NULL,
    postcode INTEGER NOT NULL,
    address_line1 VARCHAR(50) NOT NULL,
    address_line2 VARCHAR(50) NOT NULL,
    tfn VARCHAR(50) UNIQUE NOT NULL,
    super_balance float8 NOT NULL
);

INSERT INTO `employee` 
VALUES (
        'a',
        'Koichi',
        'Z',
        'Sugi',
        'male',
        5000.00,
        '1994/01/23',
        'koichi.sugi01@gmail.com',
        '0432437760',
        'VIC',
        3145,
        'ABC sesame street',
        'Cauldfiled East',
        '123414',
        1000.00
    );

INSERT INTO `employee` (
        id,
        first_name,
        middle_name,
        last_name,
        gender,
        salary,
        dob,
        email,
        phone,
        `state`,
        postcode,
        address_line1,
        address_line2,
        tfn,
        super_balance
    )
VALUES (
        'b',
        'Masa',
        'Z',
        'Kohama',
        'male',
        5000.00,
        '1994/01/23',
        'koichi.sugi02@gmail.com',
        0432437760,
        'VIC',
        3145,
        'ABC sesame street',
        'Cauldfiled East',
        '023414',
        1000.00
    );