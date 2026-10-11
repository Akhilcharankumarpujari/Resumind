import {Link} from "react-router";
import ScoreCircle from "~/components/ScoreCircle";

const ResumeCard = ({ resume: { id, companyName, jobTitle, feedback } }: { resume: Resume }) => {
    return (
        <Link to={`/resume/${id}`} className="resume-card animate-in fade-in duration-1000 shadow-md border border-gray-100 hover:shadow-lg transition-shadow duration-300">
            <div className="resume-card-header w-full">
                <div className="flex flex-col gap-1 items-start text-left w-full pr-4">
                    {companyName ? (
                        <h2 className="!text-black font-bold text-xl tracking-tight truncate w-full">{companyName}</h2>
                    ) : (
                        <h2 className="!text-black font-bold text-xl tracking-tight">Resume</h2>
                    )}
                    {jobTitle && <h3 className="text-sm break-words text-gray-400 font-medium line-clamp-1">{jobTitle}</h3>}
                </div>
                <div className="flex-shrink-0">
                    <ScoreCircle score={feedback?.overallScore || 0} />
                </div>
            </div>

            {}
            <div className="gradient-border mt-4 flex-1 flex flex-col justify-between p-4 bg-gray-50/50 rounded-xl border border-gray-100/50">
                <div className="space-y-2.5">
                    <div className="h-2.5 bg-gray-200/80 rounded w-5/6" />
                    <div className="h-2 bg-gray-200/60 rounded w-full" />
                    <div className="h-2 bg-gray-200/60 rounded w-4/5" />
                    <div className="h-2 bg-gray-200/40 rounded w-2/3" />
                </div>
                <div className="flex justify-between items-center text-xs text-gray-400 mt-6 pt-3 border-t border-gray-200/30">
                    <span>Analysis Report</span>
                    <span className="text-indigo-500 font-semibold flex items-center gap-1 group-hover:translate-x-0.5 transition-transform">
                      Open Report
                      <svg className="w-3.5 h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M9 5l7 7-7 7"></path>
                      </svg>
                    </span>
                </div>
            </div>
        </Link>
    )
}
export default ResumeCard

