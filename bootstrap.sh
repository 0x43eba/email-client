#!/bin/bash
# create_db.sh
# This script creates an SQLite database file named emails.db
# and populates it with pre-defined business-oriented example emails,
# now including a "to" column to match the Go code.

DB_FILE="emails.db"

# Remove the existing database if it exists.
if [ -f "$DB_FILE" ]; then
  echo "Removing existing database file: $DB_FILE"
  rm "$DB_FILE"
fi

echo "Creating SQLite database file: $DB_FILE"

sqlite3 "$DB_FILE" <<'EOF'
-- Drop the emails table if it already exists.
DROP TABLE IF EXISTS emails;

-- Create the emails table, now including a "to" column.
CREATE TABLE emails (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    folder TEXT NOT NULL,
    sender TEXT NOT NULL,
    subject TEXT NOT NULL,
    date TEXT NOT NULL,
    body TEXT NOT NULL,
    "to" TEXT
);

-- Insert sample emails into the Inbox. We can store the recipient in "to".
INSERT INTO emails (folder, sender, subject, date, body, "to") VALUES 
('Inbox', 'Michael Johnson', 'Quarterly Financial Update', '2023-08-15',
 'Dear Team,

Please find attached the updated financial reports for Q2. Let''s discuss the new budget proposals in our upcoming meeting.

Regards,
Michael Johnson',
 'me@company.com'),

('Inbox', 'Karen Smith', 'Client Onboarding Process', '2023-08-16',
 'Hello,

As part of our new client acquisition strategy, we are implementing an updated onboarding process. Your input during the review meeting is appreciated.

Best,
Karen Smith',
 'me@company.com');

-- Insert a sample email into the Sent folder, typically "Me" is the sender and "to" is the recipient(s).
INSERT INTO emails (folder, sender, subject, date, body, "to") VALUES 
('Sent', 'Me', 'Project Status Update', '2023-08-17',
 'Dear Team,

I am writing to update you on the current status of our major project. We are on track to meet all key deliverables. Please review the attached project plan and provide your feedback.

Regards,
[Your Name]',
 'client@external.com');

-- Insert a sample email into the Archive folder.
INSERT INTO emails (folder, sender, subject, date, body, "to") VALUES
('Archive', 'Robert Lee', 'Vendor Contract Renewal', '2023-08-14',
 'Dear Procurement,

The vendor contract is due for renewal next month. Kindly review the new terms and conditions provided and let me know your thoughts as soon as possible.

Sincerely,
Robert Lee',
 'procurement@company.com');

-- Insert a sample email into the Deleted folder.
INSERT INTO emails (folder, sender, subject, date, body, "to") VALUES 
('Deleted', 'Susan Adams', 'Team Building Activity Reminder', '2023-08-10',
 'Hello,

This is a reminder for the team building event scheduled for next week. Please confirm your attendance at your earliest convenience.

Thank you,
Susan Adams',
 'staff@company.com');

-- Insert a sample email into the Spam folder.
INSERT INTO emails (folder, sender, subject, date, body, "to") VALUES
('Spam', 'Unknown', 'You Won a Free Vacation!', '2023-08-12',
 'Congratulations!

You have been selected for a free vacation package. Please reply with your personal details to claim your prize.',
 '');
EOF

echo "Database $DB_FILE created and populated successfully."
