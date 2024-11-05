INSERT INTO categories (category) VALUES 
('General'),
('Travel'),
('Health'),
('Weather'),
('Art'),
('Fitness'),
('Books');

INSERT INTO users (email, first_name, last_name, username, password, age, gender, created_at) VALUES 
('admin@kood.tech', 'Admin', 'Admin', 'admin', '$2a$10$TZb5NJ8c.rS10oS1eDRpe.gIcuSNqCc.WODYIL2XGDIUDxPDazLYS', '34', 'Male', '2021-01-01 12:00:00 UTC'),
('moderator@kood.tech', 'Moderator', 'Moderator', 'moderator', '$2a$10$TZb5NJ8c.rS10oS1eDRpe.gIcuSNqCc.WODYIL2XGDIUDxPDazLYS','34', 'Male', '2021-01-01 12:00:00 UTC'),
('jane.doe@example.com', 'Jane', 'Doe', 'janedoe', '$2a$10$2bv7L29kab.Xr8s/i3fsZ.Asbj082x5YAlInFu08rJMGpd1yKzg62', '23', 'Female','2022-02-01 12:00:00 UTC'),
('bob.smith@example.com', 'Bob', 'Smith', 'bobsmith', '$2a$10$Wvn5k8w8.8R0P37EnP7VM.kCAUqhnTcUAjWKLqP4XegdyeBdyPcPW', '32', 'Male', '2022-03-01 12:00:00 UTC'),
('alice.johnson@example.com', 'Alice', 'Johnson', 'alicejohnson', '$2a$10$cg7X1OxxR/2R7EeQHbH0..nu5qPWvRt9EYZF3vwJunSdbxC0pbF2e', '21', 'Female', '2022-04-01 12:00:00 UTC'),
('chris.brown@example.com', 'Chris', 'Brown', 'chrisbrown', '$2a$10$MIaZSWvjsrgVMPFyaP1jX.I/V2IBM.3OOhgMChqdlRV1mKm.Hkpgy', '32', 'Male', '2022-05-01 12:00:00 UTC'),
('emily.davis@example.com', 'Emily', 'Davis', 'emilydavis', '$2a$10$1JrCALZP1gJPx5u4kw6qe.M0Fvp/lVxssohYecQ2qwQVgTSpV38A2', '24', 'Female', '2022-06-01 12:00:00 UTC'),
('vvv@vv.vv', 'Viktoriia', 'Av', 'vikvi', '$2a$10$8FhIPRwFrltybDJG7sEe0.HQgo96aEB8V6Ys1Sh/MmQ.k8DvT5ga2', '21', 'Female', '2022-06-02 12:00:00 UTC');

INSERT INTO posts (user_ID, title, content, created_at) VALUES 
(1, 'My first post', 'This is my first post!', '2021-01-01 12:00:00 UTC'),
(2, 'Amazing city Lviv', 'Just came from Lviv. It was amazing trip.', '2022-02-03 12:05:00 UTC'),
(3, 'Which vitamins are better to take', 'Need a list what better to take for sleep fixing.', '2023-01-01 12:00:00 UTC'),
(4, 'Weather in Estonia', 'What the fuck is going on?', '2022-04-01 12:00:00 UTC'),
(5, 'My latest painting', 'I just painted a Mona Lisa!', '2022-08-01 12:00:00 UTC'),
(6, 'My workout routine', 'Sharing my workout routine for getting fit!', '2022-09-01 12:00:00 UTC'),
(7, 'I''ve got a new book', 'Do you have any thoughts about "A Time to Kill" by John Grisham?', '2022-10-01 12:00:00 UTC'),
(8, 'Tallinn', 'Best buildings are in Tallinn. I draw few!', '2022-10-02 12:00:00 UTC');


INSERT INTO comments (post_ID, user_ID, content, created_at) VALUES
(1, 8, 'Great post!', '2023-01-01 12:00:00 UTC'),
(2, 7, 'Agree!', '2021-02-01 15:00:00 UTC'),
(3, 6, 'Melatonin!!!','2022-02-01 12:00:00 UTC'),
(4, 5, 'Welcome to Estonia))))))','2022-03-01 12:00:00 UTC'),
(5, 4, 'WOW!', '2022-04-01 12:00:00 UTC'),
(6, 3, 'It''s amazing', '2022-05-01 12:00:00 UTC'),
(7, 2, 'Haven''t read this book yet but heard a lot of interesting things about it.', '2022-06-01 12:00:00 UTC'),
(8, 1, 'More than agree.', '2022-07-01 12:00:00 UTC');

INSERT INTO post_categories (post_ID, category_ID) VALUES
(1, 1),
(2, 2),
(3, 3),
(4, 4),
(5, 5),
(6, 6),
(6, 3),
(7, 1),
(7, 7),
(8, 2),
(8, 5);
