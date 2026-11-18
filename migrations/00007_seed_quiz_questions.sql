-- +goose Up

-- 1. Seed broader bookstore tags
INSERT INTO tags (name) VALUES 
    ('Fantasy'),
    ('Science Fiction'),
    ('Romance'),
    ('Historical Fiction'),
    ('Horror'),
    ('Literary Fiction'),
    ('Young Adult'),
    ('Dystopian'),
    ('Contemporary')
ON CONFLICT (name) DO NOTHING;

-- 2. Insert Questions
INSERT INTO quiz_questions (prompt, position) VALUES
    ('What kind of world do you want to step into?', 1),
    ('What emotional experience are you craving?', 2),
    ('Who would you love following as the main character?', 3),
    ('What kind of pacing fits your current reading mood?', 4),
    ('What kind of stakes hook you the most?', 5);

-- 3. Insert Options for each Question
-- Question 1: World
INSERT INTO quiz_options (question_id, label) VALUES
    ((SELECT id FROM quiz_questions WHERE position = 1), 'A magical realm with deep lore and mythical power'),
    ((SELECT id FROM quiz_questions WHERE position = 1), 'The gritty, shadowed corners of our real world'),
    ((SELECT id FROM quiz_questions WHERE position = 1), 'A futuristic or dystopian society shaped by technology'),
    ((SELECT id FROM quiz_questions WHERE position = 1), 'A past era rich with historical detail and drama');

-- Question 2: Emotional Experience
INSERT INTO quiz_options (question_id, label) VALUES
    ((SELECT id FROM quiz_questions WHERE position = 2), 'A racing heartbeat with nonstop adrenaline and danger'),
    ((SELECT id FROM quiz_questions WHERE position = 2), 'Intense curiosity: trying to solve a puzzle before the reveal'),
    ((SELECT id FROM quiz_questions WHERE position = 2), 'Deep emotional connection and butterflies'),
    ((SELECT id FROM quiz_questions WHERE position = 2), 'Psychological tension that keeps me questioning reality');

-- Question 3: Protagonist
INSERT INTO quiz_options (question_id, label) VALUES
    ((SELECT id FROM quiz_questions WHERE position = 3), 'A weary investigator or detective who never gives up'),
    ((SELECT id FROM quiz_questions WHERE position = 3), 'An elite operative or warrior fighting against impossible odds'),
    ((SELECT id FROM quiz_questions WHERE position = 3), 'Two complicated people whose lives collide unexpectedly'),
    ((SELECT id FROM quiz_questions WHERE position = 3), 'An underdog trying to expose corruption or survive a deadly conspiracy');

-- Question 4: Pacing
INSERT INTO quiz_options (question_id, label) VALUES
    ((SELECT id FROM quiz_questions WHERE position = 4), 'Breakneck speed: short chapters and cliffhangers galore'),
    ((SELECT id FROM quiz_questions WHERE position = 4), 'Atmospheric and immersive: a world I can savor and live in'),
    ((SELECT id FROM quiz_questions WHERE position = 4), 'A steady, methodical investigation where every clue matters');

-- Question 5: Stakes
INSERT INTO quiz_options (question_id, label) VALUES
    ((SELECT id FROM quiz_questions WHERE position = 5), 'Personal: family secrets, deep betrayals, or true love'),
    ((SELECT id FROM quiz_questions WHERE position = 5), 'Global or national: saving a country, planet, or society'),
    ((SELECT id FROM quiz_questions WHERE position = 5), 'Justice: putting a dangerous culprit behind bars');

-- 4. Map Option Tag Weights
-- Question 1 Options
INSERT INTO option_tag_weights (option_id, tag_id, weight) VALUES
    ((SELECT id FROM quiz_options WHERE label = 'A magical realm with deep lore and mythical power'), (SELECT id FROM tags WHERE name = 'Fantasy'), 5),
    ((SELECT id FROM quiz_options WHERE label = 'A magical realm with deep lore and mythical power'), (SELECT id FROM tags WHERE name = 'Action & Adventure'), 3),
    ((SELECT id FROM quiz_options WHERE label = 'The gritty, shadowed corners of our real world'), (SELECT id FROM tags WHERE name = 'Crime'), 5),
    ((SELECT id FROM quiz_options WHERE label = 'The gritty, shadowed corners of our real world'), (SELECT id FROM tags WHERE name = 'Mystery'), 4),
    ((SELECT id FROM quiz_options WHERE label = 'The gritty, shadowed corners of our real world'), (SELECT id FROM tags WHERE name = 'Scottish Noir'), 3),
    ((SELECT id FROM quiz_options WHERE label = 'The gritty, shadowed corners of our real world'), (SELECT id FROM tags WHERE name = 'Police Procedural'), 3),
    ((SELECT id FROM quiz_options WHERE label = 'A futuristic or dystopian society shaped by technology'), (SELECT id FROM tags WHERE name = 'Science Fiction'), 5),
    ((SELECT id FROM quiz_options WHERE label = 'A futuristic or dystopian society shaped by technology'), (SELECT id FROM tags WHERE name = 'Dystopian'), 4),
    ((SELECT id FROM quiz_options WHERE label = 'A past era rich with historical detail and drama'), (SELECT id FROM tags WHERE name = 'Historical Fiction'), 5),
    ((SELECT id FROM quiz_options WHERE label = 'A past era rich with historical detail and drama'), (SELECT id FROM tags WHERE name = 'Literary Fiction'), 3);

-- Question 2 Options
INSERT INTO option_tag_weights (option_id, tag_id, weight) VALUES
    ((SELECT id FROM quiz_options WHERE label = 'A racing heartbeat with nonstop adrenaline and danger'), (SELECT id FROM tags WHERE name = 'Thriller'), 5),
    ((SELECT id FROM quiz_options WHERE label = 'A racing heartbeat with nonstop adrenaline and danger'), (SELECT id FROM tags WHERE name = 'Action & Adventure'), 4),
    ((SELECT id FROM quiz_options WHERE label = 'A racing heartbeat with nonstop adrenaline and danger'), (SELECT id FROM tags WHERE name = 'Military Fiction'), 4),
    ((SELECT id FROM quiz_options WHERE label = 'Intense curiosity: trying to solve a puzzle before the reveal'), (SELECT id FROM tags WHERE name = 'Mystery'), 5),
    ((SELECT id FROM quiz_options WHERE label = 'Intense curiosity: trying to solve a puzzle before the reveal'), (SELECT id FROM tags WHERE name = 'Crime'), 4),
    ((SELECT id FROM quiz_options WHERE label = 'Intense curiosity: trying to solve a puzzle before the reveal'), (SELECT id FROM tags WHERE name = 'Police Procedural'), 3),
    ((SELECT id FROM quiz_options WHERE label = 'Deep emotional connection and butterflies'), (SELECT id FROM tags WHERE name = 'Romance'), 5),
    ((SELECT id FROM quiz_options WHERE label = 'Deep emotional connection and butterflies'), (SELECT id FROM tags WHERE name = 'Contemporary'), 4),
    ((SELECT id FROM quiz_options WHERE label = 'Psychological tension that keeps me questioning reality'), (SELECT id FROM tags WHERE name = 'Psychological Thriller'), 5),
    ((SELECT id FROM quiz_options WHERE label = 'Psychological tension that keeps me questioning reality'), (SELECT id FROM tags WHERE name = 'Horror'), 4);

-- Question 3 Options
INSERT INTO option_tag_weights (option_id, tag_id, weight) VALUES
    ((SELECT id FROM quiz_options WHERE label = 'A weary investigator or detective who never gives up'), (SELECT id FROM tags WHERE name = 'Police Procedural'), 5),
    ((SELECT id FROM quiz_options WHERE label = 'A weary investigator or detective who never gives up'), (SELECT id FROM tags WHERE name = 'Crime'), 4),
    ((SELECT id FROM quiz_options WHERE label = 'A weary investigator or detective who never gives up'), (SELECT id FROM tags WHERE name = 'Scottish Noir'), 4),
    ((SELECT id FROM quiz_options WHERE label = 'An elite operative or warrior fighting against impossible odds'), (SELECT id FROM tags WHERE name = 'Military Fiction'), 5),
    ((SELECT id FROM quiz_options WHERE label = 'An elite operative or warrior fighting against impossible odds'), (SELECT id FROM tags WHERE name = 'Action & Adventure'), 5),
    ((SELECT id FROM quiz_options WHERE label = 'Two complicated people whose lives collide unexpectedly'), (SELECT id FROM tags WHERE name = 'Romance'), 5),
    ((SELECT id FROM quiz_options WHERE label = 'Two complicated people whose lives collide unexpectedly'), (SELECT id FROM tags WHERE name = 'Contemporary'), 4),
    ((SELECT id FROM quiz_options WHERE label = 'An underdog trying to expose corruption or survive a deadly conspiracy'), (SELECT id FROM tags WHERE name = 'Financial Thriller'), 5),
    ((SELECT id FROM quiz_options WHERE label = 'An underdog trying to expose corruption or survive a deadly conspiracy'), (SELECT id FROM tags WHERE name = 'Thriller'), 4);

-- Question 4 Options
INSERT INTO option_tag_weights (option_id, tag_id, weight) VALUES
    ((SELECT id FROM quiz_options WHERE label = 'Breakneck speed: short chapters and cliffhangers galore'), (SELECT id FROM tags WHERE name = 'Thriller'), 5),
    ((SELECT id FROM quiz_options WHERE label = 'Breakneck speed: short chapters and cliffhangers galore'), (SELECT id FROM tags WHERE name = 'Action & Adventure'), 4),
    ((SELECT id FROM quiz_options WHERE label = 'Atmospheric and immersive: a world I can savor and live in'), (SELECT id FROM tags WHERE name = 'Fantasy'), 4),
    ((SELECT id FROM quiz_options WHERE label = 'Atmospheric and immersive: a world I can savor and live in'), (SELECT id FROM tags WHERE name = 'Historical Fiction'), 5),
    ((SELECT id FROM quiz_options WHERE label = 'Atmospheric and immersive: a world I can savor and live in'), (SELECT id FROM tags WHERE name = 'Scottish Noir'), 4),
    ((SELECT id FROM quiz_options WHERE label = 'A steady, methodical investigation where every clue matters'), (SELECT id FROM tags WHERE name = 'Mystery'), 5),
    ((SELECT id FROM quiz_options WHERE label = 'A steady, methodical investigation where every clue matters'), (SELECT id FROM tags WHERE name = 'Police Procedural'), 5),
    ((SELECT id FROM quiz_options WHERE label = 'A steady, methodical investigation where every clue matters'), (SELECT id FROM tags WHERE name = 'Crime'), 4);

-- Question 5 Options
INSERT INTO option_tag_weights (option_id, tag_id, weight) VALUES
    ((SELECT id FROM quiz_options WHERE label = 'Personal: family secrets, deep betrayals, or true love'), (SELECT id FROM tags WHERE name = 'Romance'), 4),
    ((SELECT id FROM quiz_options WHERE label = 'Personal: family secrets, deep betrayals, or true love'), (SELECT id FROM tags WHERE name = 'Psychological Thriller'), 4),
    ((SELECT id FROM quiz_options WHERE label = 'Global or national: saving a country, planet, or society'), (SELECT id FROM tags WHERE name = 'Action & Adventure'), 5),
    ((SELECT id FROM quiz_options WHERE label = 'Global or national: saving a country, planet, or society'), (SELECT id FROM tags WHERE name = 'Military Fiction'), 5),
    ((SELECT id FROM quiz_options WHERE label = 'Global or national: saving a country, planet, or society'), (SELECT id FROM tags WHERE name = 'Science Fiction'), 5),
    ((SELECT id FROM quiz_options WHERE label = 'Global or national: saving a country, planet, or society'), (SELECT id FROM tags WHERE name = 'Dystopian'), 4),
    ((SELECT id FROM quiz_options WHERE label = 'Justice: putting a dangerous culprit behind bars'), (SELECT id FROM tags WHERE name = 'Crime'), 5),
    ((SELECT id FROM quiz_options WHERE label = 'Justice: putting a dangerous culprit behind bars'), (SELECT id FROM tags WHERE name = 'Police Procedural'), 5),
    ((SELECT id FROM quiz_options WHERE label = 'Justice: putting a dangerous culprit behind bars'), (SELECT id FROM tags WHERE name = 'Mystery'), 4);

-- +goose Down
DELETE FROM quiz_questions WHERE position IN (1, 2, 3, 4, 5);
DELETE FROM tags WHERE name IN (
    'Fantasy', 'Science Fiction', 'Romance', 'Historical Fiction', 
    'Horror', 'Literary Fiction', 'Young Adult', 'Dystopian', 'Contemporary'
);