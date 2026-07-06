#!/usr/bin/env python3
"""Generate frontend/public/audio/<slug>.mp3 for every vocab string Spoony
speaks, using the offline Piper TTS engine. Re-run this whenever new
words/sentences are added to the trainer or alphabet content.

Setup (one time):
    brew install ffmpeg
    pipx install piper-tts
    mkdir -p scripts/tts/voice && cd scripts/tts/voice
    curl -sLO https://huggingface.co/rhasspy/piper-voices/resolve/main/en/en_US/lessac/medium/en_US-lessac-medium.onnx
    curl -sLO https://huggingface.co/rhasspy/piper-voices/resolve/main/en/en_US/lessac/medium/en_US-lessac-medium.onnx.json

Usage:
    Add any new word/sentence to strings.txt (one per line, exact text as
    passed to speakWord() in the app), then:
        python3 scripts/tts/generate_audio.py
    New files land in frontend/public/audio/ — commit them along with the
    content change. The script skips slugs that already have an mp3, so it's
    safe to re-run after adding just a few new lines.

The slug function here MUST stay in sync with slugify() in
frontend/src/services/tts.js — that's how the app maps a spoken string to
its audio file at runtime.
"""
import re
import subprocess
import sys
import unicodedata
from pathlib import Path

HERE = Path(__file__).parent
REPO_ROOT = HERE.parent.parent
STRINGS_FILE = HERE / "strings.txt"
VOICE_MODEL = HERE / "voice" / "en_US-lessac-medium.onnx"
OUT_DIR = REPO_ROOT / "frontend" / "public" / "audio"
PIPER_BIN = str(Path.home() / ".local" / "bin" / "piper")


def slugify(s: str) -> str:
    s = s.lower()
    s = unicodedata.normalize("NFD", s)
    s = "".join(c for c in s if unicodedata.category(c) != "Mn")
    s = re.sub(r"[^a-z0-9]+", "-", s)
    return s.strip("-")


def main():
    if not VOICE_MODEL.exists():
        sys.exit(f"Voice model not found at {VOICE_MODEL} — see setup instructions in this file's docstring.")

    lines = [l.rstrip("\n") for l in STRINGS_FILE.open(encoding="utf-8")]
    lines = [l for l in lines if l.strip()]

    slugs = {}
    collisions = []
    for line in lines:
        slug = slugify(line)
        if slug in slugs and slugs[slug] != line:
            collisions.append((slug, slugs[slug], line))
        slugs[slug] = line

    if collisions:
        print("COLLISIONS FOUND — two different strings would overwrite the same file:")
        for slug, a, b in collisions:
            print(f"  {slug!r}: {a!r} vs {b!r}")
        sys.exit(1)

    print(f"{len(lines)} strings, {len(slugs)} unique slugs, no collisions.")

    OUT_DIR.mkdir(parents=True, exist_ok=True)
    ok, generated, failed = 0, 0, []
    for i, line in enumerate(lines, 1):
        slug = slugify(line)
        wav_path = OUT_DIR / f"{slug}.wav"
        mp3_path = OUT_DIR / f"{slug}.mp3"
        if mp3_path.exists():
            ok += 1
            continue
        try:
            subprocess.run(
                [PIPER_BIN, "-m", str(VOICE_MODEL), "--length-scale", "1.15", "-f", str(wav_path)],
                input=line.encode("utf-8"),
                check=True,
                capture_output=True,
            )
            subprocess.run(
                ["ffmpeg", "-y", "-loglevel", "error", "-i", str(wav_path),
                 "-codec:a", "libmp3lame", "-qscale:a", "4", str(mp3_path)],
                check=True, capture_output=True,
            )
            wav_path.unlink()
            ok += 1
            generated += 1
        except subprocess.CalledProcessError as e:
            failed.append((line, e.stderr.decode("utf-8", "replace")[:200]))
        if i % 25 == 0:
            print(f"  {i}/{len(lines)}...")

    print(f"Done: {ok} ok ({generated} newly generated), {len(failed)} failed.")
    for line, err in failed:
        print(f"FAILED {line!r}: {err}")


if __name__ == "__main__":
    main()
