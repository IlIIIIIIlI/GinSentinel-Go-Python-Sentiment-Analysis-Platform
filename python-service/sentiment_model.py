import re
import asyncio
import logging
import random
from collections import Counter
from typing import Dict, List, Any

# Configure logging
logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s - %(name)s - %(levelname)s - %(message)s'
)
logger = logging.getLogger(__name__)

# Define sentiment lexicons (sample lists - would be much larger in a real system)
POSITIVE_WORDS = {
    'en': ['good', 'great', 'excellent', 'amazing', 'awesome', 'wonderful', 'fantastic',
           'happy', 'joy', 'love', 'like', 'positive', 'beautiful', 'nice', 'perfect'],
    'zh': ['好', '优秀', '卓越', '精彩', '优质', '美好', '出色', '高兴', '快乐',
           '喜欢', '爱', '积极', '美丽', '漂亮', '完美']
}

NEGATIVE_WORDS = {
    'en': ['bad', 'terrible', 'awful', 'horrible', 'poor', 'negative', 'sad', 'angry',
           'hate', 'dislike', 'disappointing', 'worst', 'failure', 'problem', 'difficult'],
    'zh': ['坏', '糟糕', '差', '可怕', '劣质', '消极', '悲伤', '愤怒', '讨厌',
           '不喜欢', '令人失望', '最差', '失败', '问题', '困难']
}

NEGATION_WORDS = {
    'en': ['not', 'no', 'never', "don't", "doesn't", "didn't", "wasn't", "weren't",
           "isn't", "aren't", "haven't", "hasn't", "won't", "wouldn't", "couldn't", "shouldn't"],
    'zh': ['不', '没', '没有', '不是', '非', '莫', '勿', '未', '别']
}

async def analyze_text(text: str, language: str = 'en') -> Dict[str, Any]:
    """
    Analyze the sentiment of the given text.

    Args:
        text: The text to analyze
        language: Language code (default: 'en' for English)

    Returns:
        Dictionary containing sentiment analysis results
    """
    logger.info(f"Analyzing text in {language}: {text[:50]}...")

    # In a real implementation, this might be a compute-intensive operation
    # so we'll simulate some processing time
    await asyncio.sleep(random.uniform(0.1, 0.5))

    # Normalize language code
    language = language.lower()[:2]  # Just use first two characters (e.g., "en-US" -> "en")

    # Default to English if the language is not supported
    if language not in POSITIVE_WORDS:
        logger.warning(f"Language {language} not supported, defaulting to English")
        language = 'en'

    # Preprocess the text
    processed_text = preprocess_text(text, language)

    # Extract tokens (words)
    tokens = tokenize(processed_text, language)

    # Calculate sentiment scores
    pos_score, neg_score = calculate_sentiment_scores(tokens, language)

    # Determine overall sentiment
    if pos_score > neg_score:
        sentiment = "positive"
        score = pos_score / (pos_score + neg_score) if (pos_score + neg_score) > 0 else 0.5
    elif neg_score > pos_score:
        sentiment = "negative"
        score = -neg_score / (pos_score + neg_score) if (pos_score + neg_score) > 0 else -0.5
    else:
        sentiment = "neutral"
        score = 0.0

    # Ensure score is in [-1, 1] range
    score = max(min(score, 1.0), -1.0)

    # Calculate confidence scores (normalized to sum to 1)
    confidence_scores = calculate_confidence_scores(pos_score, neg_score)

    # Extract keywords
    keywords = extract_keywords(tokens, language, limit=5)

    result = {
        "sentiment": sentiment,
        "score": score,
        "confidence_scores": confidence_scores,
        "keywords": keywords
    }

    logger.info(f"Analysis complete: {sentiment} with score {score:.2f}")
    return result

def preprocess_text(text: str, language: str) -> str:
    """Preprocess the text for sentiment analysis."""
    # Convert to lowercase
    text = text.lower()

    # Remove URLs
    text = re.sub(r'https?://\S+|www\.\S+', '', text)

    # Remove email addresses
    text = re.sub(r'\S+@\S+', '', text)

    # Remove extra whitespace
    text = re.sub(r'\s+', ' ', text).strip()

    # For Chinese, we might want to use a different approach
    if language == 'zh':
        # Chinese-specific preprocessing could be added here
        pass

    return text

def tokenize(text: str, language: str) -> List[str]:
    """Tokenize the text into words or characters."""
    if language == 'zh':
        # For Chinese, character-level tokenization
        # In a real implementation, you would use a proper Chinese tokenizer
        tokens = list(text)
    else:
        # For other languages, simple word tokenization
        tokens = text.split()

        # Remove punctuation
        tokens = [re.sub(r'[^\w\s]', '', token) for token in tokens]
        # Remove empty tokens
        tokens = [token for token in tokens if token]

    return tokens

def calculate_sentiment_scores(tokens: List[str], language: str) -> tuple:
    """Calculate positive and negative sentiment scores based on lexicons."""
    positive_words = set(POSITIVE_WORDS.get(language, POSITIVE_WORDS['en']))
    negative_words = set(NEGATIVE_WORDS.get(language, NEGATIVE_WORDS['en']))
    negation_words = set(NEGATION_WORDS.get(language, NEGATION_WORDS['en']))

    pos_score = 0.0
    neg_score = 0.0

    # Simple negation handling
    negate = False

    for i, token in enumerate(tokens):
        # Check if this token is a negation word
        if token in negation_words:
            negate = True
            continue

        # Reset negation after a few tokens (simple window)
        if negate and i > 0 and (i - tokens.index(list(negation_words & set(tokens))[0]) if negation_words & set(tokens) else 0) > 3:
            negate = False

        # Check sentiment
        if token in positive_words:
            if negate:
                neg_score += 1.0
            else:
                pos_score += 1.0
        elif token in negative_words:
            if negate:
                pos_score += 0.5  # Negated negative is less positive than an actual positive
            else:
                neg_score += 1.0

        # Reset negation after a sentiment word
        if token in positive_words or token in negative_words:
            negate = False

    return pos_score, neg_score

def calculate_confidence_scores(pos_score: float, neg_score: float) -> Dict[str, float]:
    """Calculate confidence scores for each sentiment class."""
    total = pos_score + neg_score + 0.1  # Adding a small value to avoid division by zero

    positive_conf = pos_score / total
    negative_conf = neg_score / total
    neutral_conf = 1.0 - (positive_conf + negative_conf)

    # Ensure neutral confidence is not negative
    neutral_conf = max(0.0, neutral_conf)

    # Normalize to ensure they sum to 1
    total_conf = positive_conf + negative_conf + neutral_conf

    return {
        "positive": positive_conf / total_conf,
        "negative": negative_conf / total_conf,
        "neutral": neutral_conf / total_conf
    }

def extract_keywords(tokens: List[str], language: str, limit: int = 5) -> List[str]:
    """Extract the most important keywords from the text."""
    # Filter out common stopwords (a real implementation would use a proper stopword list)
    stopwords = {
        'en': ['the', 'a', 'an', 'and', 'or', 'but', 'is', 'are', 'was', 'were',
               'be', 'been', 'being', 'in', 'on', 'at', 'to', 'for', 'with', 'by'],
        'zh': ['的', '了', '和', '是', '在', '我', '有', '他', '这', '中', '你', '那', '要', '就', '人']
    }

    # Use the correct stopword list or default to English
    stop_set = set(stopwords.get(language, stopwords['en']))

    # Filter out stopwords and very short tokens
    filtered_tokens = [token for token in tokens if token not in stop_set and len(token) > 1]

    # Count word frequencies
    word_counts = Counter(filtered_tokens)

    # Get the most common words
    keywords = [word for word, count in word_counts.most_common(limit)]

    return keywords

# Add a test function for direct usage
if __name__ == "__main__":
    async def test():
        # Test English text
        result = await analyze_text("I love this product! It's amazing and works great.", "en")
        print("English test result:", result)

        # Test negative English text
        result = await analyze_text("This is terrible. I hate it and it doesn't work at all.", "en")
        print("Negative English test result:", result)

        # Test Chinese text
        result = await analyze_text("这个产品非常好，我很喜欢！", "zh")
        print("Chinese test result:", result)

        # Test with negation
        result = await analyze_text("This is not bad at all, I actually quite like it.", "en")
        print("Negation test result:", result)

    asyncio.run(test())