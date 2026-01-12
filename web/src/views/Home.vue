<template>
  <div class="min-h-screen">
    <!-- Hero section -->
    <div class="relative overflow-hidden">
      <div class="max-w-7xl mx-auto">
        <div class="relative z-10 pb-8 bg-transparent sm:pb-16 md:pb-20 lg:max-w-2xl lg:w-full lg:pb-28 xl:pb-32">
          <main class="mt-10 mx-auto max-w-7xl px-4 sm:mt-12 sm:px-6 md:mt-16 lg:mt-20 lg:px-8 xl:mt-28">
            <div class="sm:text-center lg:text-left">
              <h1 class="text-4xl tracking-tight font-extrabold text-slate-800 dark:text-slate-100 sm:text-5xl md:text-6xl">
                <span class="block xl:inline">XCStrings</span>
                <!-- Add a space and then the highlight -->
                <span class="block text-mint xl:inline"> Translator Studio</span>
              </h1>
              <p class="mt-3 text-base text-slate-600 dark:text-slate-300 sm:mt-5 sm:text-lg sm:max-w-xl sm:mx-auto md:mt-5 md:text-xl lg:mx-0">
                The most powerful and intuitive tool for translating XCStrings files. Upload your localization files, 
                select target languages, and run batch translations with your favorite provider. Export the updated 
                file in one click.
              </p>
              <div class="mt-5 sm:mt-8 sm:flex sm:justify-center lg:justify-start">
                <div class="rounded-md shadow">
                  <router-link 
                    to="/translate" 
                    class="w-full flex items-center justify-center px-8 py-3 border border-transparent text-base font-medium rounded-md text-midnight bg-mint hover:bg-mint/90 md:py-4 md:text-lg md:px-10"
                  >
                    Get started
                  </router-link>
                </div>
                <div class="mt-3 sm:mt-0 sm:ml-3">
                  <router-link 
                    to="/projects" 
                    class="w-full flex items-center justify-center px-8 py-3 border border-transparent text-base font-medium rounded-md text-mint bg-mint/10 hover:bg-mint/20 md:py-4 md:text-lg md:px-10"
                  >
                    View Projects
                  </router-link>
                </div>
              </div>
            </div>
          </main>
        </div>
      </div>
      <div class="lg:absolute lg:inset-y-0 lg:right-0 lg:w-1/2">
        <div class="h-56 w-full bg-gradient-to-r from-mint/10 to-blue-500/10 sm:h-72 md:h-96 lg:w-full lg:h-full dark:from-mint/5 dark:to-blue-500/5 rounded-l-3xl">
          <div class="h-full w-full flex items-center justify-center p-8">
            <div class="bg-white/10 dark:bg-midnight/50 backdrop-blur-lg border border-white/20 dark:border-white/10 rounded-2xl p-6 w-full max-w-md">
              <div class="space-y-4">
                <div class="h-4 bg-white/20 rounded w-3/4"></div>
                <div class="h-4 bg-white/20 rounded w-1/2"></div>
                <div class="h-4 bg-white/20 rounded w-5/6"></div>
                <div class="h-4 bg-white/20 rounded w-2/3"></div>
                <div class="h-4 bg-white/20 rounded w-4/5"></div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Features section -->
    <div class="py-12 bg-white/5 dark:bg-slate-900/5">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="lg:text-center">
          <h2 class="text-base text-mint font-semibold tracking-wide uppercase">Features</h2>
          <p class="mt-2 text-3xl leading-8 font-extrabold tracking-tight text-slate-800 dark:text-slate-100 sm:text-4xl">
            A better way to handle localization
          </p>
          <p class="mt-4 max-w-2xl text-xl text-slate-600 dark:text-slate-300 lg:mx-auto">
            Our platform provides everything you need to manage your localization workflow efficiently.
          </p>
        </div>

        <div class="mt-10">
          <div class="space-y-10 md:space-y-0 md:grid md:grid-cols-2 md:gap-x-8 md:gap-y-10">
            <div class="relative">
              <div class="absolute flex items-center justify-center h-12 w-12 rounded-md bg-mint/10 text-mint">
                <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 5h12M9 3v2m1.048 9.5A18.022 18.022 0 016.412 9m6.088 9h7M11 21l5-10 5 10M12.751 5C11.783 10.77 8.07 15.61 3 18.129" />
                </svg>
              </div>
              <div class="ml-16">
                <h3 class="text-lg font-medium text-slate-800 dark:text-slate-100">Multiple Translation Providers</h3>
                <p class="mt-2 text-base text-slate-600 dark:text-slate-300">
                  Choose from OpenAI, Google, DeepL, or Baidu to get the best translations for your content.
                </p>
              </div>
            </div>

            <div class="relative">
              <div class="absolute flex items-center justify-center h-12 w-12 rounded-md bg-mint/10 text-mint">
                <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 10V3L4 14h7v7l9-11h-7z" />
                </svg>
              </div>
              <div class="ml-16">
                <h3 class="text-lg font-medium text-slate-800 dark:text-slate-100">Batch Processing</h3>
                <p class="mt-2 text-base text-slate-600 dark:text-slate-300">
                  Process large localization files quickly with concurrent translation requests.
                </p>
              </div>
            </div>

            <div class="relative">
              <div class="absolute flex items-center justify-center h-12 w-12 rounded-md bg-mint/10 text-mint">
                <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
                </svg>
              </div>
              <div class="ml-16">
                <h3 class="text-lg font-medium text-slate-800 dark:text-slate-100">Project Management</h3>
                <p class="mt-2 text-base text-slate-600 dark:text-slate-300">
                  Organize your localization work with projects and track progress across teams.
                </p>
              </div>
            </div>

            <div class="relative">
              <div class="absolute flex items-center justify-center h-12 w-12 rounded-md bg-mint/10 text-mint">
                <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 7v10c0 2.21 3.582 4 8 4s8-1.79 8-4V7M4 7c0 2.21 3.582 4 8 4s8-1.79 8-4M4 7c0-2.21 3.582-4 8-4s8 1.79 8 4m0 5c0 2.21-3.582 4-8 4s-8-1.79-8-4" />
                </svg>
              </div>
              <div class="ml-16">
                <h3 class="text-lg font-medium text-slate-800 dark:text-slate-100">Export Ready Files</h3>
                <p class="mt-2 text-base text-slate-600 dark:text-slate-300">
                  Export your translated files in the exact format needed for your iOS applications.
                </p>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- CTA section -->
    <div class="bg-white dark:bg-midnight">
      <div class="max-w-7xl mx-auto py-12 px-4 sm:px-6 lg:py-16 lg:px-8 lg:flex lg:items-center lg:justify-between">
        <h2 class="text-3xl font-extrabold tracking-tight text-slate-900 dark:text-slate-100 sm:text-4xl">
          <span class="block">Ready to get started?</span>
          <span class="block text-mint">Start your first project today.</span>
        </h2>
        <div class="mt-8 flex lg:mt-0 lg:flex-shrink-0">
          <div class="inline-flex rounded-md shadow">
            <router-link 
              to="/register" 
              class="inline-flex items-center justify-center px-5 py-3 border border-transparent text-base font-medium rounded-md text-white bg-mint hover:bg-mint/90"
            >
              Get started
            </router-link>
          </div>
          <div class="ml-3 inline-flex rounded-md shadow">
            <router-link 
              to="/login" 
              class="inline-flex items-center justify-center px-5 py-3 border border-transparent text-base font-medium rounded-md text-mint bg-mint/10 hover:bg-mint/20"
            >
              Sign in
            </router-link>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
// Home component for the XCStrings Translator Studio
</script>
