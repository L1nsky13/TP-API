
(function() {
    const audioPlayer = document.getElementById('f1-audio-player');
    
    if (!audioPlayer) return;

    // État courant de la musique
    let isInitializing = true;

    // Fonction pour restaurer et jouer la musique
    function restoreAudioState() {
        const savedTime = parseFloat(localStorage.getItem('audioCurrentTime')) || 0;
        const wasPlaying = localStorage.getItem('audioWasPlaying') === 'true';
        const currentSrc = audioPlayer.querySelector('source').src;
        const savedSrc = localStorage.getItem('audioSrc');

        // Si c'est la même musique, continuer depuis là où elle s'est arrêtée
        if (savedSrc === currentSrc) {
            audioPlayer.currentTime = savedTime;
            // Si elle était en cours de lecture, continuer
            if (wasPlaying) {
                audioPlayer.play().catch(e => console.log('Lecture automatique bloquée:', e));
            }
        } else {
            // Nouvelle musique, réinitialiser
            localStorage.setItem('audioSrc', currentSrc);
            localStorage.setItem('audioCurrentTime', '0');
            localStorage.setItem('audioWasPlaying', 'true');
            // Jouer la nouvelle musique automatiquement
            audioPlayer.play().catch(e => console.log('Lecture automatique bloquée:', e));
        }
        isInitializing = false;
    }

    // Restaurer l'état au premier chargement ou en revenant en arrière
    window.addEventListener('pageshow', function() {
        restoreAudioState();
    });

    // Mettre à jour le temps de lecture en continu
    audioPlayer.addEventListener('timeupdate', function() {
        if (!isInitializing) {
            localStorage.setItem('audioCurrentTime', audioPlayer.currentTime);
        }
    });

    // Suivi de l'état de lecture
    audioPlayer.addEventListener('play', function() {
        localStorage.setItem('audioWasPlaying', 'true');
    });

    audioPlayer.addEventListener('pause', function() {
        localStorage.setItem('audioWasPlaying', 'false');
    });

    // Sauvegarder l'état avant de quitter la page
    window.addEventListener('beforeunload', function() {
        localStorage.setItem('audioCurrentTime', audioPlayer.currentTime);
        localStorage.setItem('audioWasPlaying', !audioPlayer.paused);
    });

    // Appeler la restauration initiale
    restoreAudioState();
})();
