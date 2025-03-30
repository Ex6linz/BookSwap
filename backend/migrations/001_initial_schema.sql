-- Aktywacja rozszerzenia dla generowania UUID
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Tabela użytkowników
CREATE TABLE users (
                       id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
                       name VARCHAR(100) NOT NULL,
                       email VARCHAR(255) UNIQUE NOT NULL,
                       password_hash VARCHAR(255) NOT NULL,
                       location VARCHAR(255),
                       bio TEXT,
                       avatar_url VARCHAR(255),
                       rating DECIMAL(3,2) DEFAULT 0,
                       created_at TIMESTAMP NOT NULL DEFAULT NOW(),
                       updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Tabela kategorii książek
CREATE TABLE categories (
                            id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
                            name VARCHAR(100) NOT NULL UNIQUE,
                            description TEXT
);

-- Tabela książek
CREATE TABLE books (
                       id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
                       title VARCHAR(255) NOT NULL,
                       author VARCHAR(255) NOT NULL,
                       description TEXT,
                       isbn VARCHAR(20),
                       category_id UUID REFERENCES categories(id),
                       condition VARCHAR(50) NOT NULL, -- np. "jak nowa", "dobra", "zużyta"
                       owner_id UUID NOT NULL REFERENCES users(id),
                       status VARCHAR(50) NOT NULL DEFAULT 'available', -- available, borrowed, reserved
                       created_at TIMESTAMP NOT NULL DEFAULT NOW(),
                       updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Tabela zdjęć książek
CREATE TABLE book_images (
                             id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
                             book_id UUID NOT NULL REFERENCES books(id) ON DELETE CASCADE,
                             image_url VARCHAR(255) NOT NULL,
                             created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Tabela transakcji wypożyczenia/wymiany
CREATE TABLE transactions (
                              id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
                              book_id UUID NOT NULL REFERENCES books(id),
                              lender_id UUID NOT NULL REFERENCES users(id), -- właściciel książki
                              borrower_id UUID NOT NULL REFERENCES users(id), -- wypożyczający
                              status VARCHAR(50) NOT NULL DEFAULT 'pending', -- pending, active, completed, canceled
                              transaction_type VARCHAR(50) NOT NULL, -- lending, exchange
                              start_date TIMESTAMP,
                              due_date TIMESTAMP,
                              return_date TIMESTAMP,
                              notes TEXT,
                              created_at TIMESTAMP NOT NULL DEFAULT NOW(),
                              updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Tabela ocen/opinii o użytkownikach
CREATE TABLE reviews (
                         id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
                         transaction_id UUID NOT NULL REFERENCES transactions(id),
                         reviewer_id UUID NOT NULL REFERENCES users(id),
                         reviewed_id UUID NOT NULL REFERENCES users(id),
                         rating INTEGER NOT NULL CHECK (rating BETWEEN 1 AND 5),
                         comment TEXT,
                         created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Tabela wiadomości między użytkownikami
CREATE TABLE messages (
                          id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
                          sender_id UUID NOT NULL REFERENCES users(id),
                          receiver_id UUID NOT NULL REFERENCES users(id),
                          transaction_id UUID REFERENCES transactions(id),
                          content TEXT NOT NULL,
                          read BOOLEAN NOT NULL DEFAULT FALSE,
                          created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Wstawienie podstawowych kategorii
INSERT INTO categories (name, description) VALUES
                                               ('Fantastyka', 'Książki fantasy, science fiction i pokrewne'),
                                               ('Kryminał', 'Kryminały, thrillery i powieści detektywistyczne'),
                                               ('Romans', 'Powieści romantyczne'),
                                               ('Literatura piękna', 'Klasyczna i współczesna literatura piękna'),
                                               ('Biografia', 'Biografie i autobiografie'),
                                               ('Historia', 'Książki historyczne'),
                                               ('Nauka', 'Książki naukowe i popularnonaukowe'),
                                               ('Dla dzieci', 'Literatura dziecięca'),
                                               ('Biznes', 'Książki biznesowe i o zarządzaniu'),
                                               ('Poradniki', 'Poradniki i książki self-help');

-- Indeksy dla zwiększenia wydajności
CREATE INDEX idx_books_owner ON books(owner_id);
CREATE INDEX idx_books_category ON books(category_id);
CREATE INDEX idx_books_status ON books(status);
CREATE INDEX idx_transactions_book ON transactions(book_id);
CREATE INDEX idx_transactions_lender ON transactions(lender_id);
CREATE INDEX idx_transactions_borrower ON transactions(borrower_id);
CREATE INDEX idx_messages_sender_receiver ON messages(sender_id, receiver_id);
CREATE INDEX idx_reviews_transaction ON reviews(transaction_id);
CREATE INDEX idx_reviews_reviewer ON reviews(reviewer_id);
CREATE INDEX idx_reviews_reviewed ON reviews(reviewed_id);


CREATE TABLE wishlist_items (
                                id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
                                user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
                                title VARCHAR(255) NOT NULL,
                                author VARCHAR(255),
                                category_id UUID REFERENCES categories(id),
                                created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Tabela powiadomień
CREATE TABLE notifications (
                               id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
                               user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
                               type VARCHAR(50) NOT NULL, -- np. 'transaction_request', 'message', 'review'
                               related_id UUID, -- id powiązanego obiektu (np. transaction_id)
                               content TEXT NOT NULL,
                               read BOOLEAN NOT NULL DEFAULT FALSE,
                               created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Indeksy dla nowych tabel
CREATE INDEX idx_wishlist_user ON wishlist_items(user_id);
CREATE INDEX idx_notifications_user ON notifications(user_id);
CREATE INDEX idx_notifications_read ON notifications(user_id, read);