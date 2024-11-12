package xml

const template = `<?xml version="1.0" encoding="UTF-8"?>
	<musicXML xmlns="http://www.musicxml.org/schema">
		<score-partwise>
			<work>
				<work-title>Untitled</work-title>
			</work>
			<part-list>
				<part id="P1">
					<score-instrument id="P1-I1">
						<instrument-name>Acoustic Grand Piano</instrument-name>
					</score-instrument>
				</part>
			</part-list>
			<part id="P1">
				<measure number="1">
					<note>
						<pitch>
							<step>C</step>
							<octave>4</octave>
						</pitch>
						<duration>4</duration>
						<type>quarter</type>
					</note>
				</measure>
			</part>
		</score-partwise>`
