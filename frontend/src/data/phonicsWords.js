// Word bank for the "letter combinations" (phonics) recognition game in the
// Alphabet section. Each entry is tagged with the single letter-combo it's
// meant to teach. A word may technically contain other combos too (e.g.
// "shark" contains both "sh" and "ar") — the game filters those out of the
// answer options so there's never more than one correct choice on screen.

export const ALL_COMBOS = [
  'OO', 'EE', 'EA', 'AI', 'AY', 'OA', 'OU', 'OW', 'OI', 'OY',
  'AR', 'OR', 'IGH', 'SH', 'CH', 'TH', 'PH', 'CK', 'NG', 'WH', 'QU',
]

// There's too much material to teach as one giant exercise, so it's split
// into several small, focused rounds — each covering a handful of related
// combos instead of all 21 at once.
export const COMBO_GROUPS = [
  { id: 'oo-ee-ea',       combos: ['OO', 'EE', 'EA'] },
  { id: 'ai-ay-oa',       combos: ['AI', 'AY', 'OA'] },
  { id: 'ou-ow-oi-oy',    combos: ['OU', 'OW', 'OI', 'OY'] },
  { id: 'ar-or-igh',      combos: ['AR', 'OR', 'IGH'] },
  { id: 'sh-ch-th',       combos: ['SH', 'CH', 'TH'] },
  { id: 'ph-ck-ng-wh-qu', combos: ['PH', 'CK', 'NG', 'WH', 'QU'] },
]

export const phonicsWords = [
  // OO — long [u:]
  { word_en: 'moon', word_ru: 'луна', emoji: '🌙', combo: 'OO' },
  { word_en: 'room', word_ru: 'комната', emoji: '🛏️', combo: 'OO' },
  { word_en: 'spoon', word_ru: 'ложка', emoji: '🥄', combo: 'OO' },
  { word_en: 'food', word_ru: 'еда', emoji: '🍽️', combo: 'OO' },
  { word_en: 'zoo', word_ru: 'зоопарк', emoji: '🦁', combo: 'OO' },
  { word_en: 'boot', word_ru: 'ботинок', emoji: '🥾', combo: 'OO' },
  { word_en: 'pool', word_ru: 'бассейн', emoji: '🏊', combo: 'OO' },
  { word_en: 'roof', word_ru: 'крыша', emoji: '🏠', combo: 'OO' },
  // OO — short [u]
  { word_en: 'book', word_ru: 'книга', emoji: '📖', combo: 'OO' },
  { word_en: 'look', word_ru: 'смотреть', emoji: '👀', combo: 'OO' },
  { word_en: 'cook', word_ru: 'готовить', emoji: '👨‍🍳', combo: 'OO' },
  { word_en: 'foot', word_ru: 'нога', emoji: '🦶', combo: 'OO' },
  { word_en: 'wood', word_ru: 'дерево', emoji: '🪵', combo: 'OO' },
  { word_en: 'good', word_ru: 'хороший', emoji: '👍', combo: 'OO' },

  // EE
  { word_en: 'tree', word_ru: 'дерево', emoji: '🌳', combo: 'EE' },
  { word_en: 'bee', word_ru: 'пчела', emoji: '🐝', combo: 'EE' },
  { word_en: 'green', word_ru: 'зелёный', emoji: '🟢', combo: 'EE' },
  { word_en: 'feet', word_ru: 'ступни', emoji: '🦶', combo: 'EE' },
  { word_en: 'sleep', word_ru: 'спать', emoji: '😴', combo: 'EE' },

  // EA
  { word_en: 'sea', word_ru: 'море', emoji: '🌊', combo: 'EA' },
  { word_en: 'tea', word_ru: 'чай', emoji: '🍵', combo: 'EA' },
  { word_en: 'eat', word_ru: 'есть', emoji: '🍴', combo: 'EA' },
  { word_en: 'read', word_ru: 'читать', emoji: '📖', combo: 'EA' },
  { word_en: 'leaf', word_ru: 'лист', emoji: '🍃', combo: 'EA' },
  { word_en: 'peach', word_ru: 'персик', emoji: '🍑', combo: 'EA' },

  // AI
  { word_en: 'rain', word_ru: 'дождь', emoji: '🌧️', combo: 'AI' },
  { word_en: 'train', word_ru: 'поезд', emoji: '🚂', combo: 'AI' },
  { word_en: 'snail', word_ru: 'улитка', emoji: '🐌', combo: 'AI' },
  { word_en: 'tail', word_ru: 'хвост', emoji: '🐕', combo: 'AI' },
  { word_en: 'paint', word_ru: 'краска', emoji: '🎨', combo: 'AI' },
  { word_en: 'nail', word_ru: 'гвоздь', emoji: '🔨', combo: 'AI' },

  // AY
  { word_en: 'day', word_ru: 'день', emoji: '☀️', combo: 'AY' },
  { word_en: 'play', word_ru: 'играть', emoji: '🎮', combo: 'AY' },
  { word_en: 'gray', word_ru: 'серый', emoji: '🌫️', combo: 'AY' },
  { word_en: 'say', word_ru: 'говорить', emoji: '💬', combo: 'AY' },
  { word_en: 'way', word_ru: 'путь', emoji: '🛣️', combo: 'AY' },

  // OA
  { word_en: 'boat', word_ru: 'лодка', emoji: '⛵', combo: 'OA' },
  { word_en: 'coat', word_ru: 'пальто', emoji: '🧥', combo: 'OA' },
  { word_en: 'goat', word_ru: 'коза', emoji: '🐐', combo: 'OA' },
  { word_en: 'road', word_ru: 'дорога', emoji: '🛤️', combo: 'OA' },
  { word_en: 'soap', word_ru: 'мыло', emoji: '🧼', combo: 'OA' },
  { word_en: 'toast', word_ru: 'тост', emoji: '🍞', combo: 'OA' },

  // OU/OW — [au]
  { word_en: 'house', word_ru: 'дом', emoji: '🏠', combo: 'OU' },
  { word_en: 'mouse', word_ru: 'мышь', emoji: '🐭', combo: 'OU' },
  { word_en: 'cloud', word_ru: 'облако', emoji: '☁️', combo: 'OU' },
  { word_en: 'cow', word_ru: 'корова', emoji: '🐄', combo: 'OW' },
  { word_en: 'brown', word_ru: 'коричневый', emoji: '🟤', combo: 'OW' },
  { word_en: 'flower', word_ru: 'цветок', emoji: '🌸', combo: 'OW' },

  // OW — second reading [ou]
  { word_en: 'snow', word_ru: 'снег', emoji: '❄️', combo: 'OW' },
  { word_en: 'window', word_ru: 'окно', emoji: '🪟', combo: 'OW' },
  { word_en: 'yellow', word_ru: 'жёлтый', emoji: '🟡', combo: 'OW' },
  { word_en: 'grow', word_ru: 'расти', emoji: '🌱', combo: 'OW' },

  // OI/OY
  { word_en: 'coin', word_ru: 'монета', emoji: '🪙', combo: 'OI' },
  { word_en: 'boy', word_ru: 'мальчик', emoji: '👦', combo: 'OY' },
  { word_en: 'toy', word_ru: 'игрушка', emoji: '🧸', combo: 'OY' },
  { word_en: 'oil', word_ru: 'масло', emoji: '🛢️', combo: 'OI' },
  { word_en: 'point', word_ru: 'точка', emoji: '📍', combo: 'OI' },

  // AR
  { word_en: 'car', word_ru: 'машина', emoji: '🚗', combo: 'AR' },
  { word_en: 'star', word_ru: 'звезда', emoji: '⭐', combo: 'AR' },
  { word_en: 'park', word_ru: 'парк', emoji: '🏞️', combo: 'AR' },
  { word_en: 'farm', word_ru: 'ферма', emoji: '🚜', combo: 'AR' },
  { word_en: 'shark', word_ru: 'акула', emoji: '🦈', combo: 'AR' },
  { word_en: 'arm', word_ru: 'рука', emoji: '💪', combo: 'AR' },

  // OR
  { word_en: 'fork', word_ru: 'вилка', emoji: '🍴', combo: 'OR' },
  { word_en: 'horse', word_ru: 'лошадь', emoji: '🐴', combo: 'OR' },
  { word_en: 'corn', word_ru: 'кукуруза', emoji: '🌽', combo: 'OR' },
  { word_en: 'door', word_ru: 'дверь', emoji: '🚪', combo: 'OR' },
  { word_en: 'storm', word_ru: 'шторм', emoji: '⛈️', combo: 'OR' },

  // IGH
  { word_en: 'night', word_ru: 'ночь', emoji: '🌃', combo: 'IGH' },
  { word_en: 'light', word_ru: 'свет', emoji: '💡', combo: 'IGH' },
  { word_en: 'right', word_ru: 'правый', emoji: '➡️', combo: 'IGH' },
  { word_en: 'high', word_ru: 'высокий', emoji: '🏔️', combo: 'IGH' },

  // SH
  { word_en: 'ship', word_ru: 'корабль', emoji: '🚢', combo: 'SH' },
  { word_en: 'shop', word_ru: 'магазин', emoji: '🏪', combo: 'SH' },
  { word_en: 'fish', word_ru: 'рыба', emoji: '🐟', combo: 'SH' },
  { word_en: 'shoe', word_ru: 'туфля', emoji: '👟', combo: 'SH' },
  { word_en: 'sheep', word_ru: 'овца', emoji: '🐑', combo: 'SH' },
  { word_en: 'brush', word_ru: 'щётка', emoji: '🪥', combo: 'SH' },

  // CH
  { word_en: 'chair', word_ru: 'стул', emoji: '🪑', combo: 'CH' },
  { word_en: 'cheese', word_ru: 'сыр', emoji: '🧀', combo: 'CH' },
  { word_en: 'chicken', word_ru: 'курица', emoji: '🐔', combo: 'CH' },
  { word_en: 'watch', word_ru: 'часы', emoji: '⌚', combo: 'CH' },
  { word_en: 'beach', word_ru: 'пляж', emoji: '🏖️', combo: 'CH' },
  { word_en: 'lunch', word_ru: 'обед', emoji: '🍱', combo: 'CH' },

  // TH — voiceless [θ]
  { word_en: 'three', word_ru: 'три', emoji: '3️⃣', combo: 'TH' },
  { word_en: 'mouth', word_ru: 'рот', emoji: '👄', combo: 'TH' },
  { word_en: 'teeth', word_ru: 'зубы', emoji: '🦷', combo: 'TH' },
  { word_en: 'bath', word_ru: 'ванна', emoji: '🛁', combo: 'TH' },
  { word_en: 'thumb', word_ru: 'большой палец', emoji: '👍', combo: 'TH' },

  // PH
  { word_en: 'phone', word_ru: 'телефон', emoji: '📱', combo: 'PH' },
  { word_en: 'photo', word_ru: 'фото', emoji: '📷', combo: 'PH' },
  { word_en: 'elephant', word_ru: 'слон', emoji: '🐘', combo: 'PH' },
  { word_en: 'dolphin', word_ru: 'дельфин', emoji: '🐬', combo: 'PH' },

  // CK
  { word_en: 'duck', word_ru: 'утка', emoji: '🦆', combo: 'CK' },
  { word_en: 'clock', word_ru: 'часы', emoji: '🕐', combo: 'CK' },
  { word_en: 'sock', word_ru: 'носок', emoji: '🧦', combo: 'CK' },
  { word_en: 'black', word_ru: 'чёрный', emoji: '⚫', combo: 'CK' },
  { word_en: 'truck', word_ru: 'грузовик', emoji: '🚚', combo: 'CK' },

  // NG
  { word_en: 'ring', word_ru: 'кольцо', emoji: '💍', combo: 'NG' },
  { word_en: 'king', word_ru: 'король', emoji: '🤴', combo: 'NG' },
  { word_en: 'sing', word_ru: 'петь', emoji: '🎤', combo: 'NG' },
  { word_en: 'wing', word_ru: 'крыло', emoji: '🪽', combo: 'NG' },
  { word_en: 'long', word_ru: 'длинный', emoji: '📏', combo: 'NG' },

  // WH
  { word_en: 'white', word_ru: 'белый', emoji: '⚪', combo: 'WH' },
  { word_en: 'whale', word_ru: 'кит', emoji: '🐋', combo: 'WH' },
  { word_en: 'wheel', word_ru: 'колесо', emoji: '🛞', combo: 'WH' },
  { word_en: 'what', word_ru: 'что', emoji: '❓', combo: 'WH' },
  { word_en: 'where', word_ru: 'где', emoji: '📍', combo: 'WH' },

  // QU
  { word_en: 'queen', word_ru: 'королева', emoji: '👸', combo: 'QU' },
  { word_en: 'question', word_ru: 'вопрос', emoji: '🤔', combo: 'QU' },
  { word_en: 'quick', word_ru: 'быстрый', emoji: '⚡', combo: 'QU' },
]
