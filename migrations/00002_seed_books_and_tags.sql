-- +goose Up
INSERT INTO tags (name, kind)
VALUES 
    ('Crime', 'genre'),
    ('Mystery', 'genre'),
    ('Thriller', 'genre'),
    ('Action & Adventure', 'genre'),
    ('Psychological Thriller', 'genre'),
    ('Police Procedural', 'genre'),
    ('Scottish Noir', 'genre'),
    ('Financial Thriller', 'genre'),
    ('Military Fiction', 'genre')
ON CONFLICT (name) DO NOTHING;

INSERT INTO books (title, author, description, price, stock, image_url, status)
VALUES 
    (
        'Let It Bleed',
        'Ian Rankin',
        'Detective Inspector John Rebus investigates a string of seemingly unrelated tragedies—two suicides, a suspicious prison death, and a brutal murder—uncovering an embezzlement conspiracy that reaches the highest echelons of Scottish government and corporate power.',
        3000,
        5,
        'https://covers.openlibrary.org/b/isbn/9780752877198-L.jpg',
        'published'
    ),
    (
        'False Prophet',
        'Faye Kellerman',
        'When the glamorous owner of an upscale Beverly Hills health spa is brutally assaulted, LAPD Detective Peter Decker and Rina Lazarus delve into a viper''s nest of family jealousy, dark deceit, and charismatic fraud.',
        3000,
        4,
        'https://covers.openlibrary.org/b/isbn/9780060750886-L.jpg',
        'published'
    ),
    (
        'Blue Gold',
        'Clive Cussler',
        'Kurt Austin and the NUMA Special Assignments Team battle a ruthless international conspiracy bent on cornering the world''s freshwater supply, leading them on a perilous mission from the Venezuelan jungle to the Mediterranean.',
        3000,
        7,
        'https://covers.openlibrary.org/b/isbn/9780743418225-L.jpg',
        'published'
    ),
    (
        'Final Venture',
        'Michael Ridpath',
        'In Boston''s cutthroat venture capital world, young British partner Simon Ayot must fight to clear his name when his father-in-law and business partner is found shot dead, with all evidence pointing directly at him.',
        3000,
        4,
        'https://covers.openlibrary.org/b/isbn/9780451197535-L.jpg',
        'published'
    ),
    (
        'Night Sins',
        'Tami Hoag',
        'A peaceful Minnesota town is shattered when an eight-year-old boy vanishes without a trace, leaving a sadistic ransom note. SBI investigator Megan O''Malley and Police Chief Mitch Holt race to unmask an elusive predator hidden within the community.',
        3000,
        5,
        'https://covers.openlibrary.org/b/isbn/9780553564518-L.jpg',
        'published'
    ),
    (
        'Zero Option',
        'Chris Ryan',
        'SAS operative Geordie Sharp is plunged into a covert counter-terrorist mission against rogue paramilitary factions and ruthless conspirators, where every move could trigger global catastrophe and betrayal lurks at every turn.',
        3000,
        6,
        'https://covers.openlibrary.org/b/isbn/9780099247623-L.jpg',
        'published'
    ),
    (
        'The Lewis Man',
        'Peter May',
        'When a mummified body is unearthed from a peat bog on the Isle of Lewis, DNA tests link the murder to an elderly local man with advanced dementia, forcing former detective Fin Macleod to unravel a decades-old family tragedy.',
        3000,
        5,
        'https://covers.openlibrary.org/b/isbn/9781849169295-L.jpg',
        'published'
    ),
    (
        'The Chessmen',
        'Peter May',
        'Working as an estate security officer on the Isle of Lewis, Fin Macleod investigates poachers, only for a mysterious natural bog burst to expose a sunken plane containing the body of a friend who vanished seventeen years before.',
        3000,
        5,
        'https://covers.openlibrary.org/b/isbn/9780857382238-L.jpg',
        'published'
    ),
    (
        'The Blackhouse',
        'Peter May',
        'Edinburgh Detective Inspector Fin Macleod returns to his harsh childhood home on the Isle of Lewis to investigate a savage copycat killing, confronting buried memories, the island''s brutal guga hunt tradition, and deep personal trauma.',
        3000,
        8,
        'https://covers.openlibrary.org/b/isbn/9781849163859-L.jpg',
        'published'
    ),
    (
        'When the Music''s Over',
        'Peter Robinson',
        'DCI Alan Banks probes historical sexual abuse allegations against high-profile celebrities, while DI Annie Cabbot investigates the horrific murder of a vulnerable young woman abandoned in the Yorkshire countryside.',
        3000,
        5,
        'https://covers.openlibrary.org/b/isbn/9780062394781-L.jpg',
        'published'
    ),
    (
        'Guilty Wives',
        'James Patterson',
        'Four affluent best friends embark on a lavish weekend getaway in Monte Carlo, but their dream holiday turns into a living nightmare when they wake up framed for murder and locked in a high-security French penitentiary.',
        3000,
        6,
        'https://covers.openlibrary.org/b/isbn/9781455506668-L.jpg',
        'published'
    ),
    (
        'No Time for Goodbye',
        'Linwood Barclay',
        'At fourteen years old, Cynthia woke up to find that her entire family had vanished without a trace, leaving breakfast on the table. Twenty-five years later, her public plea for answers awakens a deadly conspiracy.',
        3000,
        7,
        'https://covers.openlibrary.org/b/isbn/9780553592283-L.jpg',
        'published'
    ),
    (
        'Too Close to Home',
        'Linwood Barclay',
        'In the peaceful suburb of Promise Falls, the Cutter family''s next-door neighbors are executed in cold blood in their home. Terror consumes the neighborhood as Jim Cutter realizes the killers may have hit the wrong house.',
        3000,
        5,
        'https://covers.openlibrary.org/b/isbn/9780553805567-L.jpg',
        'published'
    ),
    (
        'The Safe House',
        'Nicci French',
        'When a prominent family is massacred in their secluded home, trauma specialist Dr. Samantha Laschen offers sanctuary to the sole survivor, an eighteen-year-old girl, inadvertently placing her own daughter in the sights of a killer.',
        3000,
        4,
        'https://covers.openlibrary.org/b/isbn/9780140270365-L.jpg',
        'published'
    ),
    (
        'Dead Souls',
        'Ian Rankin',
        'Inspector John Rebus contends with ghosts from the past across a stormy Edinburgh: an old friend''s suicide, a convicted predator returning to a local housing scheme, and a cold murder inquiry that sparks widespread panic.',
        3000,
        6,
        'https://covers.openlibrary.org/b/isbn/9780752826288-L.jpg',
        'published'
    ),
    (
        'Global Strike',
        'Chris Ryan',
        'SAS warrior John Porter is deployed on an off-the-books mission to intercept extremist warlords armed with devastating bioweapons before they can orchestrate a coordinated international strike.',
        3000,
        5,
        'https://covers.openlibrary.org/b/isbn/9781444783810-L.jpg',
        'published'
    ),
    (
        'Deception on His Mind',
        'Elizabeth George',
        'While Thomas Lynley is away on his honeymoon, Detective Sergeant Barbara Havers is sent to an Essex seaside town to investigate the racially charged murder of a young Pakistani immigrant, exposing an undercurrent of community bigotry.',
        3000,
        5,
        'https://covers.openlibrary.org/b/isbn/9780553575093-L.jpg',
        'published'
    ),
    (
        'Skinner''s Round',
        'Quintin Jardine',
        'Edinburgh Assistant Chief Constable Bob Skinner is taking time off for a golf tournament in East Lothian when an American guest is discovered brutally slain, hurling Skinner into an intricate web of international intrigue and ancient grudges.',
        3000,
        4,
        'https://covers.openlibrary.org/b/isbn/9780747250418-L.jpg',
        'published'
    );

-- Tag Associations
INSERT INTO book_tags (book_id, tag_id)
SELECT b.id, t.id
FROM books b, tags t
WHERE 
    (b.title = 'Let It Bleed' AND t.name IN ('Crime', 'Mystery', 'Scottish Noir'))
    OR (b.title = 'False Prophet' AND t.name IN ('Mystery', 'Crime', 'Thriller'))
    OR (b.title = 'Blue Gold' AND t.name IN ('Action & Adventure', 'Thriller'))
    OR (b.title = 'Final Venture' AND t.name IN ('Financial Thriller', 'Thriller', 'Mystery'))
    OR (b.title = 'Night Sins' AND t.name IN ('Psychological Thriller', 'Thriller', 'Mystery'))
    OR (b.title = 'Zero Option' AND t.name IN ('Military Fiction', 'Action & Adventure', 'Thriller'))
    OR (b.title = 'The Lewis Man' AND t.name IN ('Mystery', 'Crime', 'Scottish Noir'))
    OR (b.title = 'The Chessmen' AND t.name IN ('Mystery', 'Crime', 'Scottish Noir'))
    OR (b.title = 'The Blackhouse' AND t.name IN ('Mystery', 'Crime', 'Scottish Noir'))
    OR (b.title = 'When the Music''s Over' AND t.name IN ('Police Procedural', 'Crime', 'Mystery'))
    OR (b.title = 'Guilty Wives' AND t.name IN ('Psychological Thriller', 'Thriller'))
    OR (b.title = 'No Time for Goodbye' AND t.name IN ('Psychological Thriller', 'Mystery', 'Thriller'))
    OR (b.title = 'Too Close to Home' AND t.name IN ('Thriller', 'Mystery', 'Crime'))
    OR (b.title = 'The Safe House' AND t.name IN ('Psychological Thriller', 'Mystery', 'Thriller'))
    OR (b.title = 'Dead Souls' AND t.name IN ('Crime', 'Mystery', 'Scottish Noir'))
    OR (b.title = 'Global Strike' AND t.name IN ('Military Fiction', 'Action & Adventure', 'Thriller'))
    OR (b.title = 'Deception on His Mind' AND t.name IN ('Police Procedural', 'Mystery', 'Crime'))
    OR (b.title = 'Skinner''s Round' AND t.name IN ('Police Procedural', 'Crime', 'Scottish Noir'));

-- +goose Down
DELETE FROM book_tags WHERE book_id IN (
    SELECT id FROM books WHERE title IN (
        'Let It Bleed', 'False Prophet', 'Blue Gold', 'Final Venture', 'Night Sins',
        'Zero Option', 'The Lewis Man', 'The Chessmen', 'The Blackhouse', 'When the Music''s Over',
        'Guilty Wives', 'No Time for Goodbye', 'Too Close to Home', 'The Safe House',
        'Dead Souls', 'Global Strike', 'Deception on His Mind', 'Skinner''s Round'
    )
);

DELETE FROM books WHERE title IN (
    'Let It Bleed', 'False Prophet', 'Blue Gold', 'Final Venture', 'Night Sins',
    'Zero Option', 'The Lewis Man', 'The Chessmen', 'The Blackhouse', 'When the Music''s Over',
    'Guilty Wives', 'No Time for Goodbye', 'Too Close to Home', 'The Safe House',
    'Dead Souls', 'Global Strike', 'Deception on His Mind', 'Skinner''s Round'
);

DELETE FROM tags WHERE name IN (
    'Crime', 'Mystery', 'Thriller', 'Action & Adventure', 'Psychological Thriller',
    'Police Procedural', 'Scottish Noir', 'Financial Thriller', 'Military Fiction'
);
