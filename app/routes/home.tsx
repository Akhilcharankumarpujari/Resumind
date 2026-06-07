import type { Route } from "./+types/home";
import Navbar from "~/components/Navbar";
import ResumeCard from "~/components/ResumeCard";
import { useAuthStore } from "~/lib/auth";
import { Link, useNavigate } from "react-router";
import { useEffect, useState } from "react";

export function meta({}: Route.MetaArgs) {
  return [
    { title: "Resumind - AI Resume Analyzer & ATS Optimization" },
    { name: "description", content: "Optimize your resume for applicant tracking systems with instant, interactive AI feedback powered by Groq Llama 3.3." },
  ];
}

interface ResumeItem {
  id: string;
  companyName: string;
  jobTitle: string;
  feedback: {
    overallScore: number;
  };
  createdAt: string;
}

export default function Home() {
  const { isAuthenticated, isLoading } = useAuthStore();
  const navigate = useNavigate();
  const [resumes, setResumes] = useState<ResumeItem[]>([]);
  const [loadingResumes, setLoadingResumes] = useState(false);


  useEffect(() => {
    if (!isAuthenticated) return;

    const loadResumes = async () => {
      setLoadingResumes(true);
      try {
        const res = await fetch("/api/resumes", { credentials: "include" });
        if (res.ok) {
          const data = await res.json();
          setResumes(data || []);
        }
      } catch (err) {
        console.error("Failed to load resumes:", err);
      } finally {
        setLoadingResumes(false);
      }
    };

    loadResumes();
  }, [isAuthenticated]);

  if (isLoading) {
    return (
      <main className="bg-gradient min-h-screen flex flex-col items-center justify-center pt-0">
        <div className="flex flex-col items-center gap-4">
          <img src="/images/resume-scan-2.gif" className="w-[180px] h-auto object-contain mix-blend-multiply" alt="Loading Resumind..." />
          <p className="text-gray-500 font-medium animate-pulse">Initializing Resumind...</p>
        </div>
      </main>
    );
  }

  if (!isAuthenticated) {
    return <LandingPage />;
  }


  return (
    <main className="bg-[url('/images/bg-main.svg')] bg-cover min-h-screen pt-4 pb-20">
      <Navbar />

      <section className="main-section px-4 max-w-7xl mx-auto">
        <div className="page-heading py-12 max-w-3xl mx-auto text-center flex flex-col items-center gap-6">
          <h1 className="text-4xl md:text-5xl lg:text-6xl font-extrabold tracking-tight text-gradient leading-tight">
            Your Resume Dashboard
          </h1>
          {!loadingResumes && resumes?.length === 0 ? (
            <h2 className="text-lg md:text-xl text-gray-500 font-normal">
              No resumes analyzed yet. Upload your first resume to get instant ATS scores and expert AI advice.
            </h2>
          ) : (
            <h2 className="text-lg md:text-xl text-gray-500 font-normal">
              Track your applications and review AI-powered ratings for your submissions.
            </h2>
          )}
        </div>

        {}
        {loadingResumes && (
          <div className="flex flex-col items-center justify-center py-12">
            <img src="/images/resume-scan-2.gif" className="w-[160px] h-auto mix-blend-multiply" alt="Loading resumes..." />
            <p className="text-gray-500 mt-4 font-medium">Loading your resumes...</p>
          </div>
        )}

        {}
        {!loadingResumes && resumes.length > 0 && (
          <div className="w-full mt-6">
            <div className="flex justify-between items-center mb-6">
              <h3 className="text-xl font-bold text-gray-800">Recent Analyses ({resumes.length})</h3>
              <Link to="/upload" className="primary-button w-fit text-sm py-2 px-5 font-medium shadow-md hover:scale-[1.02] transition-transform duration-200">
                + New Analysis
              </Link>
            </div>
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6 justify-center">
              {resumes.map((resume) => (
                <ResumeCard key={resume.id} resume={resume as any} />
              ))}
            </div>
          </div>
        )}

        {}
        {!loadingResumes && resumes?.length === 0 && (
          <div className="flex flex-col items-center justify-center mt-6 p-10 bg-white/60 backdrop-blur-md border border-gray-200/50 rounded-2xl max-w-md mx-auto text-center shadow-lg">
            <div className="bg-indigo-50 p-4 rounded-full mb-4">
              <svg className="w-10 h-10 text-indigo-500" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"></path>
              </svg>
            </div>
            <h4 className="text-lg font-bold text-gray-800 mb-2">Get Started Now</h4>
            <p className="text-sm text-gray-500 mb-6">
              Align your resume with a specific job description and get detailed ratings for ATS readiness, formatting, content, and skill alignment.
            </p>
            <Link to="/upload" className="primary-button text-base py-3 px-8 font-semibold shadow-md hover:scale-[1.02] transition-all">
              Upload Your Resume
            </Link>
          </div>
        )}
      </section>
    </main>
  );
}


function LandingPage() {
  const [activeTab, setActiveTab] = useState<"ats" | "tone" | "skills" | "structure">("ats");

  const interactiveMockData = {
    ats: {
      score: 84,
      title: "ATS Compatibility Score",
      description: "Checks how easily automated systems can read and parse your resume structure.",
      tips: [
        { type: "good", tip: "Standard section titles found", explanation: "You used 'Experience' and 'Education' which ATS parsers recognize instantly." },
        { type: "good", tip: "Clean tableless layout", explanation: "No tables or complex column grids found. Perfect for basic linear parsing." },
        { type: "improve", tip: "Avoid text boxes", explanation: "You have 2 text boxes containing contact info. Most ATS systems bypass text boxes entirely, leaving your contact details blank." },
        { type: "improve", tip: "Uncommon date format", explanation: "Using '2020 - Current' is less parsable than '08/2020 - Present'." }
      ]
    },
    tone: {
      score: 72,
      title: "Tone & Style Evaluation",
      description: "Measures the impact, action-orientation, and readability of your writing.",
      tips: [
        { type: "good", tip: "Strong action verbs in bullets", explanation: "Began 85% of your bullet points with active verbs like 'Designed', 'Orchestrated', and 'Optimized'." },
        { type: "improve", tip: "Passive voice usage", explanation: "Detected passive phrasing: 'Was responsible for managing client relationships'. Change to: 'Managed high-value client relationships'." },
        { type: "improve", tip: "Vague results", explanation: "Phrases like 'Improved server performance' lack metrics. Revise to: 'Boosted server speed by 35% through query optimization'." }
      ]
    },
    skills: {
      score: 65,
      title: "Skills Gap & Keyword Alignment",
      description: "Compares keywords in your resume directly against the target job description.",
      tips: [
        { type: "good", tip: "Core technology match", explanation: "Found exact matches for 'Go', 'React', and 'MySQL' as required by the Job Description." },
        { type: "improve", tip: "Missing high-priority keywords", explanation: "The JD heavily mentions 'CI/CD' and 'Docker', which do not appear in your resume." },
        { type: "improve", tip: "Unstructured skills section", explanation: "Your skills are listed in a generic block. Categorize them into 'Languages', 'Frameworks', and 'Tools' for readability." }
      ]
    },
    structure: {
      score: 90,
      title: "Formatting & Structure Rating",
      description: "Validates the visual hierarchy, margins, and document flow.",
      tips: [
        { type: "good", tip: "Consistent font styling", explanation: "Font sizes and weights are consistently used across headings, subheadings, and body copy." },
        { type: "good", tip: "Healthy margin space", explanation: "Maintained a professional 0.75-inch margin around the document, avoiding a cluttered look." },
        { type: "improve", tip: "Excessive length", explanation: "Your experience spans 3 pages, but has only 4 years of total experience. Condense into a punchy 1-page document." }
      ]
    }
  };

  const currentTab = interactiveMockData[activeTab];

  return (
    <div className="bg-gradient min-h-screen text-gray-900 selection:bg-indigo-200 selection:text-indigo-900 pb-20">
      {}
      <div className="pt-4 px-4 max-w-7xl mx-auto">
        <Navbar />
      </div>

      {}
      <section className="relative pt-16 pb-20 px-6 text-center max-w-5xl mx-auto overflow-hidden">
        {}
        <div className="absolute top-10 left-1/2 -translate-x-1/2 w-[500px] h-[500px] bg-indigo-200/40 rounded-full blur-[120px] pointer-events-none -z-10 animate-pulse" />
        <div className="absolute top-20 left-1/4 w-[300px] h-[300px] bg-indigo-200/20 rounded-full blur-[100px] pointer-events-none -z-10" />
        <div className="animate-fade-in flex flex-col items-center gap-6">
          <h1 className="text-4xl sm:text-5xl md:text-7xl font-extrabold tracking-tight leading-tight max-w-4xl text-gradient">
            Pass the ATS. <br />Land the Interview.
          </h1>

          <p className="text-lg md:text-xl text-gray-500 max-w-2xl mt-2 leading-relaxed">
            Upload your resume, paste the target job description, and get instant, honest grading with granular tips for matching, styling, and ATS readability.
          </p>

          <div className="flex flex-col sm:flex-row gap-4 mt-6 w-full sm:w-auto justify-center">
            <Link to="/upload" className="primary-button text-lg font-semibold py-4 px-10 shadow-xl shadow-indigo-200 hover:scale-[1.03] transition-all text-center flex items-center justify-center gap-2">
              Optimize Resume Free
              <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M13 7l5 5m0 0l-5 5m5-5H6"></path>
              </svg>
            </Link>
            <a href="#how-it-works" className="bg-white/80 backdrop-blur-md border border-gray-200 rounded-full py-4 px-8 text-base font-semibold text-gray-700 hover:bg-gray-50 transition-all text-center hover:scale-[1.01]">
              How it works
            </a>
          </div>
        </div>

        {}
        <div className="mt-16 gradient-border max-w-4xl mx-auto shadow-2xl relative group">
          <div className="bg-white rounded-xl overflow-hidden shadow-inner border border-gray-100 flex flex-col md:flex-row">
            {}
            <div className="w-full md:w-5/12 bg-gray-50 p-6 border-r border-gray-100 text-left flex flex-col gap-4">
              <div className="flex items-center gap-2 border-b border-gray-200 pb-3">
                <div className="w-3 h-3 rounded-full bg-red-400" />
                <div className="w-3 h-3 rounded-full bg-yellow-400" />
                <div className="w-3 h-3 rounded-full bg-green-400" />
                <span className="text-xs text-gray-400 font-mono ml-2">resume_draft.pdf</span>
              </div>
              <div className="h-4 bg-gray-200 rounded w-2/3" />
              <div className="h-3 bg-gray-200 rounded w-1/2" />
              <div className="space-y-2 mt-4">
                <div className="h-2 bg-gray-200 rounded w-full" />
                <div className="h-2 bg-gray-200 rounded w-full" />
                <div className="h-2 bg-gray-200 rounded w-5/6" />
              </div>
              <div className="space-y-2 mt-4">
                <div className="h-3 bg-gray-200 rounded w-1/3" />
                <div className="h-2 bg-gray-200 rounded w-full" />
                <div className="h-2 bg-gray-200 rounded w-4/5" />
              </div>
            </div>

            {}
            <div className="w-full md:w-7/12 p-8 text-left flex flex-col justify-between">
              <div>
                <div className="flex justify-between items-center mb-4">
                  <span className="text-xs font-bold uppercase tracking-wider text-indigo-500">Instant AI Report</span>
                  <span className="bg-indigo-50 text-indigo-700 text-xs px-2.5 py-1 rounded-full font-bold">Llama 3.3 LPU</span>
                </div>
                <div className="flex items-center gap-4 mb-6">
                  <div className="w-16 h-16 rounded-full bg-indigo-50 border-4 border-indigo-500 flex items-center justify-center text-xl font-black text-indigo-600">
                    84
                  </div>
                  <div>
                    <h4 className="font-bold text-gray-800 text-lg">Senior Software Engineer</h4>
                    <p className="text-sm text-gray-400">Google Inc. / ATS Scorecard</p>
                  </div>
                </div>

                <div className="space-y-3">
                  <div className="flex items-start gap-2.5 text-sm">
                    <span className="text-green-500 font-bold mt-0.5">✓</span>
                    <div>
                      <p className="font-semibold text-gray-700">Excellent Keyword Density</p>
                      <p className="text-gray-400 text-xs">Successfully matched key terms 'Go', 'MySQL' and 'Fiber'.</p>
                    </div>
                  </div>
                  <div className="flex items-start gap-2.5 text-sm">
                    <span className="text-amber-500 font-bold mt-0.5">⚠</span>
                    <div>
                      <p className="font-semibold text-gray-700">Passive Action Verbs Found</p>
                      <p className="text-gray-400 text-xs">Replace passive responsibilities with active achievement-based phrases.</p>
                    </div>
                  </div>
                </div>
              </div>

              <div className="border-t border-gray-100 pt-4 mt-6 flex justify-between items-center text-xs text-gray-400">
                <span>Analysis time: 1.4 seconds</span>
                <span className="text-indigo-500 font-semibold cursor-pointer group-hover:translate-x-1 transition-transform">See full report →</span>
              </div>
            </div>
          </div>
        </div>
      </section>

      {}
      <section className="bg-white/40 border-y border-gray-100 py-10 px-6 text-center">
        <p className="text-xs font-semibold uppercase tracking-widest text-gray-400 mb-6">
          Optimize your resume for applications at leading companies
        </p>
        <div className="flex flex-wrap items-center justify-center gap-8 md:gap-16 opacity-50 grayscale hover:grayscale-0 transition-all duration-300">
          <span className="text-lg md:text-2xl font-black tracking-tight text-gray-800">Google</span>
          <span className="text-lg md:text-2xl font-black tracking-tight text-gray-800">Meta</span>
          <span className="text-lg md:text-2xl font-black tracking-tight text-gray-800">Microsoft</span>
          <span className="text-lg md:text-2xl font-black tracking-tight text-gray-800">Amazon</span>
          <span className="text-lg md:text-2xl font-black tracking-tight text-gray-800">Netflix</span>
        </div>
      </section>

      {}
      <section id="how-it-works" className="py-20 px-6 max-w-6xl mx-auto scroll-mt-10">
        <div className="text-center max-w-2xl mx-auto mb-16 flex flex-col gap-4">
          <h2 className="text-3xl md:text-4xl font-extrabold text-gray-900 tracking-tight">
            Three Steps to Interview Callbacks
          </h2>
          <p className="text-gray-500">
            Resumind takes the guesswork out of optimizing your resume for ATS parsers and hiring managers.
          </p>
        </div>

        <div className="grid grid-cols-1 md:grid-cols-3 gap-8 relative">
          <div className="absolute top-1/2 left-0 right-0 h-0.5 bg-gradient-to-r from-indigo-100 via-rose-100 to-indigo-100 -translate-y-1/2 hidden md:block -z-10" />

          {}
          <div className="bg-white/80 backdrop-blur-md p-8 rounded-2xl border border-gray-200/50 shadow-lg text-center flex flex-col items-center gap-4 hover:scale-[1.02] transition-transform duration-300">
            <div className="w-12 h-12 rounded-full bg-indigo-50 border border-indigo-100 flex items-center justify-center text-lg font-bold text-indigo-600 shadow-inner">
              1
            </div>
            <h3 className="text-xl font-bold text-gray-800">Sign In with Google</h3>
            <p className="text-sm text-gray-500 leading-relaxed">
              Create your account in seconds using Google Identity login. We keep your analyses saved securely in your personalized dashboard.
            </p>
          </div>

          {}
          <div className="bg-white/80 backdrop-blur-md p-8 rounded-2xl border border-gray-200/50 shadow-lg text-center flex flex-col items-center gap-4 hover:scale-[1.02] transition-transform duration-300">
            <div className="w-12 h-12 rounded-full bg-rose-50 border border-rose-100 flex items-center justify-center text-lg font-bold text-rose-600 shadow-inner">
              2
            </div>
            <h3 className="text-xl font-bold text-gray-800">Align Your Goal</h3>
            <p className="text-sm text-gray-500 leading-relaxed">
              Upload your PDF resume and paste the target job description. We extract full-text content from all pages of the document.
            </p>
          </div>

          {}
          <div className="bg-white/80 backdrop-blur-md p-8 rounded-2xl border border-gray-200/50 shadow-lg text-center flex flex-col items-center gap-4 hover:scale-[1.02] transition-transform duration-300">
            <div className="w-12 h-12 rounded-full bg-indigo-50 border border-indigo-100 flex items-center justify-center text-lg font-bold text-indigo-600 shadow-inner">
              3
            </div>
            <h3 className="text-xl font-bold text-gray-800">Analyze & Rewrite</h3>
            <p className="text-sm text-gray-500 leading-relaxed">
              Receive a granular breakdown with concrete tips to boost keyword alignment, repair parsing issues, and style your bullet points.
            </p>
          </div>
        </div>
      </section>

      {}
      <section className="py-20 px-6 bg-white/40 border-y border-gray-100">
        <div className="max-w-6xl mx-auto">
          <div className="text-center max-w-2xl mx-auto mb-16 flex flex-col gap-4">
            <h2 className="text-3xl md:text-4xl font-extrabold text-gray-900 tracking-tight">
              Interactive Preview: How We Rate Your Resume
            </h2>
            <p className="text-gray-500">
              Explore the interactive dashboard mock below to see the granularity of our assessment metrics.
            </p>
          </div>

          <div className="grid grid-cols-1 lg:grid-cols-12 gap-8 items-stretch">
            {}
            <div className="lg:col-span-4 flex flex-row lg:flex-col gap-2 overflow-x-auto lg:overflow-x-visible pb-2 lg:pb-0">
              <button
                onClick={() => setActiveTab("ats")}
                className={`flex-1 lg:flex-none text-left p-4 rounded-xl font-bold border transition-all flex items-center justify-between gap-4 cursor-pointer min-w-[150px] lg:min-w-0 ${
                  activeTab === "ats"
                    ? "bg-white border-indigo-200 shadow-md text-indigo-700"
                    : "bg-white/50 border-gray-200/50 text-gray-500 hover:bg-white/80"
                }`}
              >
                <span>ATS Parsing</span>
                <span className={`text-xs px-2 py-0.5 rounded-full font-bold ${activeTab === "ats" ? "bg-indigo-100 text-indigo-700" : "bg-gray-100 text-gray-500"}`}>84</span>
              </button>

              <button
                onClick={() => setActiveTab("tone")}
                className={`flex-1 lg:flex-none text-left p-4 rounded-xl font-bold border transition-all flex items-center justify-between gap-4 cursor-pointer min-w-[150px] lg:min-w-0 ${
                  activeTab === "tone"
                    ? "bg-white border-indigo-200 shadow-md text-indigo-700"
                    : "bg-white/50 border-gray-200/50 text-gray-500 hover:bg-white/80"
                }`}
              >
                <span>Tone & Style</span>
                <span className={`text-xs px-2 py-0.5 rounded-full font-bold ${activeTab === "tone" ? "bg-indigo-100 text-indigo-700" : "bg-gray-100 text-gray-500"}`}>72</span>
              </button>

              <button
                onClick={() => setActiveTab("skills")}
                className={`flex-1 lg:flex-none text-left p-4 rounded-xl font-bold border transition-all flex items-center justify-between gap-4 cursor-pointer min-w-[150px] lg:min-w-0 ${
                  activeTab === "skills"
                    ? "bg-white border-indigo-200 shadow-md text-indigo-700"
                    : "bg-white/50 border-gray-200/50 text-gray-500 hover:bg-white/80"
                }`}
              >
                <span>Skills Gap</span>
                <span className={`text-xs px-2 py-0.5 rounded-full font-bold ${activeTab === "skills" ? "bg-indigo-100 text-indigo-700" : "bg-gray-100 text-gray-500"}`}>65</span>
              </button>

              <button
                onClick={() => setActiveTab("structure")}
                className={`flex-1 lg:flex-none text-left p-4 rounded-xl font-bold border transition-all flex items-center justify-between gap-4 cursor-pointer min-w-[150px] lg:min-w-0 ${
                  activeTab === "structure"
                    ? "bg-white border-indigo-200 shadow-md text-indigo-700"
                    : "bg-white/50 border-gray-200/50 text-gray-500 hover:bg-white/80"
                }`}
              >
                <span>Structure & Layout</span>
                <span className={`text-xs px-2 py-0.5 rounded-full font-bold ${activeTab === "structure" ? "bg-indigo-100 text-indigo-700" : "bg-gray-100 text-gray-500"}`}>90</span>
              </button>
            </div>

            {}
            <div className="lg:col-span-8 bg-white border border-gray-200/60 rounded-2xl shadow-xl p-6 md:p-8 text-left flex flex-col gap-6 relative overflow-hidden">
              <div className="absolute top-0 right-0 w-32 h-32 bg-indigo-500/5 rounded-full blur-2xl pointer-events-none" />

              <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-4 border-b border-gray-100 pb-5">
                <div>
                  <h3 className="text-xl font-extrabold text-gray-800">{currentTab.title}</h3>
                  <p className="text-sm text-gray-400 mt-1">{currentTab.description}</p>
                </div>
                <div className="flex items-center gap-3 bg-gray-50 p-2.5 rounded-xl border border-gray-100 self-start sm:self-auto">
                  <span className="text-xs uppercase font-bold tracking-wider text-gray-400">Section Score</span>
                  <div className="w-12 h-12 rounded-lg bg-indigo-600 text-white font-black flex items-center justify-center text-lg shadow-md shadow-indigo-100">
                    {currentTab.score}
                  </div>
                </div>
              </div>

              <div className="space-y-4">
                {currentTab.tips.map((tip, idx) => (
                  <div key={idx} className="flex items-start gap-4 p-4 rounded-xl border border-gray-50 hover:bg-gray-50/50 transition-colors">
                    {tip.type === "good" ? (
                      <span className="w-6 h-6 rounded-full bg-emerald-50 text-emerald-600 flex items-center justify-center font-bold text-xs shrink-0 border border-emerald-100">✓</span>
                    ) : (
                      <span className="w-6 h-6 rounded-full bg-amber-50 text-amber-600 flex items-center justify-center font-bold text-xs shrink-0 border border-amber-100">!</span>
                    )}
                    <div>
                      <h4 className="font-bold text-gray-800 text-sm">{tip.tip}</h4>
                      <p className="text-xs text-gray-500 mt-1 leading-relaxed">{tip.explanation}</p>
                    </div>
                  </div>
                ))}
              </div>
            </div>
          </div>
        </div>
      </section>

      {}
      <section className="py-20 px-6 max-w-6xl mx-auto">
        <div className="text-center max-w-2xl mx-auto mb-16 flex flex-col gap-4">
          <h2 className="text-3xl md:text-4xl font-extrabold text-gray-900 tracking-tight">
            Features Optimized for Speed and Precision
          </h2>
          <p className="text-gray-500">
            Everything you need to bypass filters and land interviews at top tech companies.
          </p>
        </div>

        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-8">
          {}
          <div className="bg-white/70 p-6 rounded-2xl border border-gray-200/50 shadow hover:shadow-lg transition-all flex flex-col gap-4">
            <div className="w-10 h-10 rounded-xl bg-indigo-50 text-indigo-600 flex items-center justify-center font-bold">
              ⚡
            </div>
            <h3 className="text-lg font-bold text-gray-800">Ultra-fast Analysis</h3>
            <p className="text-sm text-gray-500 leading-relaxed">
              Powered by Groq's specialized inference engine. The entire multi-page PDF parsing and AI grading takes less than 3 seconds.
            </p>
          </div>

          {}
          <div className="bg-white/70 p-6 rounded-2xl border border-gray-200/50 shadow hover:shadow-lg transition-all flex flex-col gap-4">
            <div className="w-10 h-10 rounded-xl bg-indigo-50 text-indigo-600 flex items-center justify-center font-bold">
              📂
            </div>
            <h3 className="text-lg font-bold text-gray-800">Multi-page Parsing</h3>
            <p className="text-sm text-gray-500 leading-relaxed">
              Unlike other free tools that only read page 1, our Go backend loops through every page in your PDF to extract all text content accurately.
            </p>
          </div>

          {}
          <div className="bg-white/70 p-6 rounded-2xl border border-gray-200/50 shadow hover:shadow-lg transition-all flex flex-col gap-4">
            <div className="w-10 h-10 rounded-xl bg-indigo-50 text-indigo-600 flex items-center justify-center font-bold">
              🎯
            </div>
            <h3 className="text-lg font-bold text-gray-800">Job Description Alignment</h3>
            <p className="text-sm text-gray-500 leading-relaxed">
              Compare your resume against specific target jobs. We pinpoint missing keywords, tools, and technical terms instantly.
            </p>
          </div>

          {}
          <div className="bg-white/70 p-6 rounded-2xl border border-gray-200/50 shadow hover:shadow-lg transition-all flex flex-col gap-4">
            <div className="w-10 h-10 rounded-xl bg-indigo-50 text-indigo-600 flex items-center justify-center font-bold">
              🔒
            </div>
            <h3 className="text-lg font-bold text-gray-800">Secure Database Storage</h3>
            <p className="text-sm text-gray-500 leading-relaxed">
              We use MySQL GORM schemas to securely log user details and preserve all your previous upload analyses in a clean dashboard.
            </p>
          </div>

          {}
          <div className="bg-white/70 p-6 rounded-2xl border border-gray-200/50 shadow hover:shadow-lg transition-all flex flex-col gap-4">
            <div className="w-10 h-10 rounded-xl bg-indigo-50 text-indigo-600 flex items-center justify-center font-bold">
              🍪
            </div>
            <h3 className="text-lg font-bold text-gray-800">Secure JWT Cookie Auth</h3>
            <p className="text-sm text-gray-500 leading-relaxed">
              Features HTTP-only cookie authentication session tokens. Rest assured that no client-side script can access your login credentials.
            </p>
          </div>

          {}
          <div className="bg-white/70 p-6 rounded-2xl border border-gray-200/50 shadow hover:shadow-lg transition-all flex flex-col gap-4">
            <div className="w-10 h-10 rounded-xl bg-indigo-50 text-indigo-600 flex items-center justify-center font-bold">
              📝
            </div>
            <h3 className="text-lg font-bold text-gray-800">Actionable Tips</h3>
            <p className="text-sm text-gray-500 leading-relaxed">
              No generic suggestions. You receive specific headlines and explanations on exactly how to rewrite each bullet point or formatting error.
            </p>
          </div>
        </div>
      </section>


      {}
      <section className="px-6 pb-12 max-w-5xl mx-auto">
        <div className="bg-gradient-to-br from-[#8e98ff] to-[#606beb] rounded-3xl p-10 md:p-16 text-center text-white relative overflow-hidden shadow-xl shadow-indigo-100 flex flex-col items-center gap-6">
          <div className="absolute -top-12 -right-12 w-64 h-64 bg-white/10 rounded-full blur-2xl pointer-events-none" />
          <div className="absolute -bottom-12 -left-12 w-64 h-64 bg-indigo-400/20 rounded-full blur-2xl pointer-events-none" />

          <h2 className="text-3xl md:text-5xl font-extrabold tracking-tight text-white">
            Ready to Optimize Your Resume?
          </h2>
          <p className="text-indigo-100 max-w-xl text-base md:text-lg">
            Create an account in seconds, upload your resume, and start matching your credentials against your dream roles.
          </p>
          <Link to="/upload" className="bg-white text-indigo-700 font-bold text-lg rounded-full py-4 px-10 mt-4 shadow-xl hover:scale-[1.03] transition-all hover:bg-indigo-50">
            Get Started Free
          </Link>
        </div>
      </section>

      {}
      <footer className="border-t border-gray-200/50 pt-10 px-6 max-w-6xl mx-auto flex flex-col sm:flex-row justify-between items-center gap-6 text-gray-400 text-xs">
        <p>© 2026 Resumind. All rights reserved. By Akhil Charan Kumar.</p>
        <div className="flex gap-6">
          <Link to="/" className="hover:text-gray-600 transition-colors">Privacy Policy</Link>
          <Link to="/" className="hover:text-gray-600 transition-colors">Terms of Service</Link>
          <a href="mailto:akhilcharankumar001@gmail.com" className="hover:text-gray-600 transition-colors">Contact Support</a>
        </div>
      </footer>
    </div>
  );
}
