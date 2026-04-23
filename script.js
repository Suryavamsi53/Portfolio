document.addEventListener('DOMContentLoaded', () => {
    // Initialize Lucide Icons
    if (window.lucide) {
        lucide.createIcons();
    }

    // Intersection Observer for Reveal Animations
    const revealEls = document.querySelectorAll('.glass-card, .section-header');
    if (revealEls.length > 0) {
        const observerOptions = {
            threshold: 0.1,
            rootMargin: '0px 0px -50px 0px'
        };

        const observer = new IntersectionObserver((entries) => {
            entries.forEach(entry => {
                if (entry.isIntersecting) {
                    entry.target.classList.add('reveal-active');
                    observer.unobserve(entry.target);
                }
            });
        }, observerOptions);

        revealEls.forEach(el => {
            el.style.opacity = '0';
            el.style.transform = 'translateY(30px)';
            el.style.transition = 'all 0.8s cubic-bezier(0.16, 1, 0.3, 1)';
            observer.observe(el);
        });

        const style = document.createElement('style');
        style.innerHTML = `.reveal-active { opacity: 1 !important; transform: translateY(0) !important; }`;
        document.head.appendChild(style);
    }

    // Form Handling
    const contactForm = document.getElementById('contact-form');
    if (contactForm) {
        contactForm.addEventListener('submit', async (e) => {
            e.preventDefault();
            const submitBtn = contactForm.querySelector('.submit-btn');
            const originalText = submitBtn.innerText;
            
            const formData = {
                name: document.getElementById('name').value,
                email: document.getElementById('email').value,
                phone: document.getElementById('phone').value,
                company: document.getElementById('company').value,
                region: document.getElementById('region').value,
                country: document.getElementById('country').value,
                message: document.getElementById('message').value
            };

            submitBtn.innerText = 'Transmitting...';
            submitBtn.disabled = true;

            try {
                const response = await fetch('api/contact', {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify(formData)
                });

                if (response.ok) {
                    submitBtn.innerText = 'Transmission Successful';
                    submitBtn.style.background = '#00d4ff';
                    contactForm.reset();
                } else {
                    throw new Error('Server returned error status');
                }
            } catch (err) {
                console.error('Transmission Error:', err);
                submitBtn.innerText = 'Error: System Offline';
                submitBtn.style.background = '#ff3b3b';
            }

            setTimeout(() => {
                submitBtn.innerText = originalText;
                submitBtn.style.background = '';
                submitBtn.disabled = false;
            }, 3000);
        });
    }

    // Smooth Navigation highlighting
    const sections = document.querySelectorAll('section');
    const navLinks = document.querySelectorAll('.nav-links a');

    if (sections.length > 0 && navLinks.length > 0) {
        window.addEventListener('scroll', () => {
            let current = '';
            sections.forEach(section => {
                const sectionTop = section.offsetTop;
                if (window.pageYOffset >= sectionTop - 150) {
                    current = section.getAttribute('id');
                }
            });

            navLinks.forEach(link => {
                link.classList.remove('active');
                if (link.getAttribute('href').includes(current)) {
                    link.classList.add('active');
                }
            });
        });

        const navStyle = document.createElement('style');
        navStyle.innerHTML = `.nav-links a.active { color: var(--accent-secondary) !important; }`;
        document.head.appendChild(navStyle);
    }
});
